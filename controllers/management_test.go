package controllers

import (
	"context"
	"net/http"
	"testing"

	"wishy/models"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHandleManagementDispatchesHTTPMethods(t *testing.T) {
	r := require.New(t)
	t.Chdir("..")

	request := events.APIGatewayProxyRequest{
		Path:       "/prod/manage",
		HTTPMethod: http.MethodGet,
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage: "prod",
		},
	}
	response, err := HandleManagement(context.Background(), request, nil)
	r.NoError(err)
	r.Equal(http.StatusOK, response.StatusCode)
	r.Contains(response.Body, "Sign in to Wishy")
	r.Contains(response.Body, `action="/prod/manage"`)

	request.HTTPMethod = http.MethodDelete
	response, err = HandleManagement(context.Background(), request, nil)
	r.NoError(err)
	r.Equal(http.StatusMethodNotAllowed, response.StatusCode)
	r.Equal("GET, POST", response.Headers["Allow"])
}

func TestManagerRedirectPreservesStagePath(t *testing.T) {
	r := require.New(t)
	request := events.APIGatewayProxyRequest{
		Path: "/manage",
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage:      "prod",
			DomainName: "abc123.execute-api.eu-west-1.amazonaws.com",
			Path:       "/prod/manage",
		},
	}
	response, err := redirectManager("session-token", request)
	r.NoError(err)
	r.Equal(http.StatusSeeOther, response.StatusCode)
	r.Equal("/prod/manage", response.Headers["Location"])
	r.Contains(response.Headers["Set-Cookie"], "Path=/;")

	logoutResponse, err := redirectManagerLogout(request)
	r.NoError(err)
	r.Equal("/prod/manage", logoutResponse.Headers["Location"])
	r.Contains(logoutResponse.Headers["Set-Cookie"], "Path=/;")

	customDomainRequest := events.APIGatewayProxyRequest{
		Path: "/manage",
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage:      "prod",
			DomainName: "wishy.example.com",
			Path:       "/prod/manage",
		},
	}
	r.Equal("/prod/manage", managementPath(customDomainRequest))
	r.Equal("/manage", NormalizedRequestPath(customDomainRequest))

	stagePrefixedRequest := events.APIGatewayProxyRequest{
		Path: "/prod/manage",
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage: "prod",
		},
	}
	r.Equal("/manage", NormalizedRequestPath(stagePrefixedRequest))
	r.Equal("/prod/manage", managementPath(stagePrefixedRequest))
}

func TestManagerPasswordVerification(t *testing.T) {
	password := "generated-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	if !verifyManagerPassword(string(hash), password) {
		t.Fatal("verifyManagerPassword() rejected the correct password")
	}
	if verifyManagerPassword(string(hash), password+"wrong") {
		t.Fatal("verifyManagerPassword() accepted an incorrect password")
	}
	if verifyManagerPassword(password, password) {
		t.Fatal("verifyManagerPassword() accepted a plaintext password")
	}
}

func TestRenderHtmlManagementPage(t *testing.T) {
	r := require.New(t)
	t.Chdir("..")
	templates.InitHtmlTpls()

	category := models.Category{BaseCategory: models.BaseCategory{Name: "<script>category</script>"}}
	wish := models.Wish{
		BaseWish: models.BaseWish{Name: "<script>wish</script>", Preference: models.Low},
		Category: category.BaseCategory,
		Active:   true,
	}
	page := managerPageData{
		Authenticated: true,
		User:          models.User{Username: "manager"},
		Categories:    []models.Category{category},
		Wishes:        []models.Wish{wish},
		Error:         "<script>alert(1)</script>",
		Path:          "/prod/manage",
	}
	response, err := renderManagerPage(page, http.StatusOK, nil)
	r.NoError(err)
	r.NotContains(response.Body, "<script>alert(1)</script>")
	r.NotContains(response.Body, "<script>category</script>")
	r.NotContains(response.Body, "<script>wish</script>")
	r.Contains(response.Body, "&lt;script&gt;")
	r.Contains(response.Body, `<html lang="en">`)
	r.Contains(response.Body, "Categories")
	r.Contains(response.Body, "Wish Management")
	r.Contains(response.Body, `action="/prod/manage"`)
	r.NotContains(response.Body, "Gestione wishes")
	r.NotContains(response.Body, "Accedi a Wishy")
}

func TestGetUserWishesErrors(t *testing.T) {
	r := require.New(t)
	response, err := getUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{}}, nil)
	r.Nil(response)
	r.Equal("400", err.Error())

	response, err = getUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"uid": ""}}, nil)
	r.Nil(response)
	r.Equal("400", err.Error())

	response, err = getUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"uid": "aasdf"}}, nil)
	r.Nil(response)
	r.Equal("400", err.Error())
}

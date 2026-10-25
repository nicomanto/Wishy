package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
	"wishy/controllers"
	"wishy/models"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

var mockWishesList models.UserWishes
var mockFriendlyError models.FriendlyError

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

// tests setup
func setup() {
	// pre load html page
	templates.InitHtmlTpls()
	mockWishesList = models.UserWishes{
		Username: "username",
		Wishes: []models.WishByCategory{
			{
				Cat: "cat-1",
				Wishes: []models.BaseWish{
					{
						Name:       "wish-cat-1-name-1",
						Link:       "wish-cat-1-link-1",
						Preference: models.Low,
					},
					{
						Name:       "wish-cat-1-name-2",
						Link:       "wish-cat-1-link-2",
						Preference: models.Medium,
					},
					{
						Name:       "wish-cat-1-name-3",
						Link:       "wish-cat-1-link-3",
						Preference: models.High,
					},
				},
			},
			{
				Cat: "cat-2",
				Wishes: []models.BaseWish{
					{
						Name:       "wish-cat-2-name-1",
						Link:       "wish-cat-2-link-1",
						Preference: models.High,
					},
					{
						Name:       "wish-cat-2-name-2",
						Link:       "wish-cat-2-link-2",
						Preference: models.High,
					},
					{
						Name:       "wish-cat-2-name-3",
						Link:       "wish-cat-2-link-3",
						Preference: models.High,
					},
				}},
		},
	}
	mockFriendlyError = models.FriendlyErrorInit("400")
}

func TestRenderHtmlWishListPage(t *testing.T) {
	r := require.New(t)
	// load html page
	var responseBody strings.Builder
	err := templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&responseBody, mockWishesList)
	r.NoError(err)
}

func TestRenderHtmlErrorPage(t *testing.T) {
	r := require.New(t)
	// load html page
	var responseBody strings.Builder
	err := templates.HtmlTpls[templates.ErrorPageHtmlTemplateType].Execute(&responseBody, mockFriendlyError)
	r.NoError(err)
}

func TestRenderHtmlManagementPage(t *testing.T) {
	r := require.New(t)
	category := models.Category{BaseCategory: models.BaseCategory{Name: "<script>category</script>"}}
	wish := models.Wish{
		BaseWish: models.BaseWish{Name: "<script>wish</script>", Preference: models.Low},
		Category: category.BaseCategory,
		Active:   true,
	}
	page := controllers.ManagerPageData{
		Authenticated: true,
		User:          models.User{Username: "manager"},
		Categories:    []models.Category{category},
		Wishes:        []models.Wish{wish},
		Error:         "<script>alert(1)</script>",
		Path:          "/prod/manage",
	}
	var responseBody strings.Builder
	err := templates.HtmlTpls[templates.ManagementHtmlTemplateType].Execute(&responseBody, page)
	r.NoError(err)
	r.NotContains(responseBody.String(), "<script>alert(1)</script>")
	r.NotContains(responseBody.String(), "<script>category</script>")
	r.NotContains(responseBody.String(), "<script>wish</script>")
	r.Contains(responseBody.String(), "&lt;script&gt;")
	r.Contains(responseBody.String(), `<html lang="en">`)
	r.Contains(responseBody.String(), "Categories")
	r.Contains(responseBody.String(), "Wish Management")
	r.Contains(responseBody.String(), `action="/prod/manage"`)
	r.NotContains(responseBody.String(), "Gestione wishes")
	r.NotContains(responseBody.String(), "Accedi a Wishy")
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
	r.Equal(303, response.StatusCode)
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
	r.Equal("/prod/manage", managerRedirectPath(customDomainRequest))
	r.Equal("/manage", normalizedRequestPath(customDomainRequest))

	stagePrefixedRequest := events.APIGatewayProxyRequest{
		Path: "/prod/manage",
		RequestContext: events.APIGatewayProxyRequestContext{
			Stage: "prod",
		},
	}
	r.Equal("/manage", normalizedRequestPath(stagePrefixedRequest))
	r.Equal("/prod/manage", managerRedirectPath(stagePrefixedRequest))
}

func TestGetUserWishesErrors(t *testing.T) {
	r := require.New(t)
	resp, err := controllers.GetUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{}}, nil)
	r.Nil(resp)
	r.Equal("400", err.Error())
	resp, err = controllers.GetUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"uid": ""}}, nil)
	r.Nil(resp)
	r.Equal("400", err.Error())
	resp, err = controllers.GetUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"uid": "aasdf"}}, nil)
	r.Nil(resp)
	r.Equal("400", err.Error())
}

func TestGenerateWishlistPDF(t *testing.T) {
	r := require.New(t)
	b, err := mockWishesList.GenerateWishlistPDF()
	r.NoError(err)
	r.NotNil(b)
}

func TestSetRecent(t *testing.T) {
	r := require.New(t)
	wish := models.BaseWish{
		Ts: time.Now().AddDate(0, -2, 0),
	}
	r.False(wish.IsRecent)
	wish.SetRecent(3)
	r.True(wish.IsRecent)
	wish.Ts = time.Now().AddDate(0, -4, 0)
	wish.SetRecent(3)
	r.False(wish.IsRecent)
}

package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"wishy/common"
	"wishy/controllers"
	"wishy/models"
	"wishy/services/mongodb"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Endpoint string

const (
	WishesEndpoint    Endpoint = "/wishes"
	WishesPDFEndpoint Endpoint = WishesEndpoint + "/pdf"
	ManageEndpoint    Endpoint = "/manage"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	client := mongodb.Connect(os.Getenv("MONGODB_URI"))
	db := client.Database(os.Getenv("DB_NAME"))

	// pre load html page
	templates.InitHtmlTpls()

	endpoint := Endpoint(request.Path)
	// Execute different logic based on the endpoint
	switch endpoint {
	case ManageEndpoint:
		return handleManagement(ctx, request, db)
	case WishesEndpoint:
		var responseBody strings.Builder
		wishes, err := controllers.GetUserWishes(ctx, request, db)
		if err != nil {
			friendlyError := models.FriendlyErrorInit(err.Error())
			// load html error page
			err = templates.HtmlTpls[templates.ErrorPageHtmlTemplateType].Execute(&responseBody, friendlyError)
			if err != nil {
				logrus.Errorln(err)
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Headers:    map[string]string{"Content-Type": "application/json"},
				}, fmt.Errorf("get wishes: %v", err)
			}
			response, e := common.HTMLResponse(responseBody.String())
			return *response, e
		}
		// load html wishlist page
		err = templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&responseBody, wishes)
		if err != nil {
			logrus.Errorln(err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, fmt.Errorf("get wishes: %v", err)
		}
		response, e := common.HTMLResponse(responseBody.String())
		return *response, e
	case WishesPDFEndpoint:
		var responseBody strings.Builder
		wishes, err := controllers.GetUserWishes(ctx, request, db)
		if err != nil {
			friendlyError := models.FriendlyErrorInit(err.Error())
			// load html error page
			err = templates.HtmlTpls[templates.ErrorPageHtmlTemplateType].Execute(&responseBody, friendlyError)
			if err != nil {
				logrus.Errorln(err)
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Headers:    map[string]string{"Content-Type": "application/json"},
				}, fmt.Errorf("get wishes: %v", err)
			}
			response, e := common.HTMLResponse(responseBody.String())
			return *response, e
		}
		// Generate PDF
		pdfBytes, err := wishes.GenerateWishlistPDF()
		if err != nil {
			friendlyError := models.FriendlyErrorInit(err.Error())
			// load html error page
			err = templates.HtmlTpls[templates.ErrorPageHtmlTemplateType].Execute(&responseBody, friendlyError)
			if err != nil {
				logrus.Errorln(err)
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusInternalServerError,
					Headers:    map[string]string{"Content-Type": "application/json"},
				}, fmt.Errorf("get wishes: %v", err)
			}
			response, e := common.HTMLResponse(responseBody.String())
			return *response, e
		}
		// Convert PDF bytes to base64 (required for AWS API Gateway)
		pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
		response, e := common.PDFResponse(pdfBase64, wishes.Username+"-wishlist.pdf", true)
		return *response, e
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, fmt.Errorf("not found")
	}
}

func handleManagement(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == http.MethodGet {
		user, err := managerUserFromRequest(ctx, request, db)
		if err != nil {
			if managerStatus(err) != http.StatusUnauthorized {
				logrus.Errorln(err)
				return renderManagerPage(controllers.ManagerPageData{Error: "An error occurred. Please try again later."}, http.StatusInternalServerError, nil)
			}
			headers := map[string]string{"Set-Cookie": expiredManagerCookie()}
			return renderManagerPage(controllers.ManagerPageData{}, http.StatusOK, headers)
		}
		page, err := controllers.LoadManagerPage(ctx, user, db)
		if err != nil {
			logrus.Errorln(err)
			return renderManagerPage(controllers.ManagerPageData{Error: "Unable to load data."}, http.StatusInternalServerError, nil)
		}
		return renderManagerPage(page, http.StatusOK, nil)
	}
	if request.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed, Headers: map[string]string{"Allow": "GET, POST"}}, nil
	}

	form, err := managementForm(request)
	if err != nil {
		return renderManagerPage(controllers.ManagerPageData{Error: "Invalid request."}, http.StatusBadRequest, nil)
	}
	action := form.Get("action")
	switch action {
	case "login":
		token, err := controllers.LoginManager(ctx, form.Get("username"), form.Get("password"), db)
		if err != nil {
			if errors.Is(err, controllers.ErrInvalidCredentials) {
				return renderManagerPage(controllers.ManagerPageData{Error: "Incorrect username or password."}, http.StatusUnauthorized, nil)
			}
			logrus.Errorln(err)
			return renderManagerPage(controllers.ManagerPageData{Error: "An error occurred. Please try again later."}, http.StatusInternalServerError, nil)
		}
		return redirectManager(token, request)
	case "logout":
		return redirectManagerLogout(request)
	}

	user, err := managerUserFromRequest(ctx, request, db)
	if err != nil {
		if managerStatus(err) != http.StatusUnauthorized {
			logrus.Errorln(err)
			return renderManagerPage(controllers.ManagerPageData{Error: "An error occurred. Please try again later."}, http.StatusInternalServerError, nil)
		}
		return renderManagerPage(controllers.ManagerPageData{Error: "Please sign in again."}, http.StatusUnauthorized, map[string]string{"Set-Cookie": expiredManagerCookie()})
	}
	if err := controllers.ManageData(ctx, action, form, user, db); err != nil {
		status := managerStatus(err)
		if status >= http.StatusInternalServerError {
			logrus.Errorln(err)
		}
		page, loadErr := controllers.LoadManagerPage(ctx, user, db)
		if loadErr != nil {
			logrus.Errorln(loadErr)
			return renderManagerPage(controllers.ManagerPageData{Error: "Unable to load data."}, http.StatusInternalServerError, nil)
		}
		if status >= http.StatusInternalServerError {
			page.Error = "An error occurred. Please try again later."
		} else {
			page.Error = err.Error()
		}
		return renderManagerPage(page, status, nil)
	}
	return redirectManager("", request)
}

func managementForm(request events.APIGatewayProxyRequest) (url.Values, error) {
	body := request.Body
	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, err
		}
		body = string(decoded)
	}
	return url.ParseQuery(body)
}

func managerUserFromRequest(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (models.User, error) {
	var cookieHeader string
	for key, value := range request.Headers {
		if strings.EqualFold(key, "Cookie") {
			cookieHeader = value
			break
		}
	}
	if cookieHeader == "" {
		for key, values := range request.MultiValueHeaders {
			if strings.EqualFold(key, "Cookie") && len(values) > 0 {
				cookieHeader = values[0]
				break
			}
		}
	}
	for _, part := range strings.Split(cookieHeader, ";") {
		name, value, found := strings.Cut(strings.TrimSpace(part), "=")
		if found && name == "wishy_manager" {
			return controllers.ManagerFromSession(ctx, value, db)
		}
	}
	return models.User{}, fmt.Errorf("%d", http.StatusUnauthorized)
}

func managerStatus(err error) int {
	var managerErr *controllers.ManagerError
	if errors.As(err, &managerErr) {
		return managerErr.StatusCode
	}
	status, parseErr := strconv.Atoi(err.Error())
	if parseErr == nil && status >= 400 && status <= 599 {
		return status
	}
	return http.StatusInternalServerError
}

func renderManagerPage(page controllers.ManagerPageData, status int, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var body strings.Builder
	if err := templates.HtmlTpls[templates.ManagementHtmlTemplateType].Execute(&body, page); err != nil {
		logrus.Errorln(err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: map[string]string{"Content-Type": "text/plain"}}, fmt.Errorf("render management page: %w", err)
	}
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Type"] = "text/html; charset=utf-8"
	headers["Cache-Control"] = "no-store"
	return events.APIGatewayProxyResponse{StatusCode: status, Headers: headers, Body: body.String()}, nil
}

func managerRedirectPath(request events.APIGatewayProxyRequest) string {
	stage := request.RequestContext.Stage
	if stage != "" && stage != "$default" && strings.Contains(request.RequestContext.DomainName, ".execute-api.") {
		return "/" + stage + "/manage"
	}
	if strings.HasSuffix(request.RequestContext.Path, "/manage") {
		return request.RequestContext.Path
	}
	if strings.HasSuffix(request.Path, "/manage") {
		return request.Path
	}
	return "/manage"
}

func redirectManager(token string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Location":      managerRedirectPath(request),
		"Cache-Control": "no-store",
	}
	if token != "" {
		headers["Set-Cookie"] = managerCookie(token)
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusSeeOther, Headers: headers}, nil
}

func redirectManagerLogout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusSeeOther,
		Headers: map[string]string{
			"Location":      managerRedirectPath(request),
			"Set-Cookie":    expiredManagerCookie(),
			"Cache-Control": "no-store",
		},
	}, nil
}

func managerCookie(token string) string {
	return "wishy_manager=" + token + "; Path=/; Max-Age=" + strconv.Itoa(int(12*time.Hour.Seconds())) + "; HttpOnly; Secure; SameSite=Strict"
}

func expiredManagerCookie() string {
	return "wishy_manager=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"
}

func main() {
	lambda.Start(handler)
}

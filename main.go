package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"wishy/common"
	"wishy/controllers"
	"wishy/models"
	"wishy/mongodb"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	client := mongodb.Connect(os.Getenv("MONGODB_URI"))
	db := client.Database(os.Getenv("DB_NAME"))

	// pre load html page
	templates.InitHtmlTpls()

	endpoint := request.Path
	// Execute different logic based on the endpoint
	switch endpoint {
	/*case "/categories":
	response, err := controllers.GetCategeories(ctx, request, db)
	if err != nil {
		errInt, e := strconv.Atoi(err.Error())
		if e != nil {
			return events.APIGatewayProxyResponse{}, fmt.Errorf("error: %v - %v", err, e)
		}
		return events.APIGatewayProxyResponse{
			StatusCode: errInt,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, fmt.Errorf("get categories: %v", err)
	}
	return *response, nil*/
	case "/wishes":
		var responseBody strings.Builder
		wishes, err := controllers.GetWishes(ctx, request, db)
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
	case "/wishes/pdf":
		/*var responseBody strings.Builder
		wishes, err := controllers.GetWishes(ctx, request, db)
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
		}*/
		// Open the file
		file, err := os.Open("./templates/sample.pdf")
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		defer file.Close()

		// Read the file content as bytes
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		// Convert PDF bytes to base64 (required for AWS API Gateway)
		//pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
		pdfBase64 := base64.StdEncoding.EncodeToString(data)
		response, e := common.PDFResponse(pdfBase64, "wishlist.pdf", true)
		return *response, e
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, fmt.Errorf("not found")
	}
}

func main() {
	lambda.Start(handler)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"wishy/controllers"
	"wishy/services/mongodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	endpoint := Endpoint(controllers.NormalizedRequestPath(request))
	// Route each endpoint to its controller.
	switch endpoint {
	case ManageEndpoint:
		return controllers.HandleManagement(ctx, request, db)
	case WishesEndpoint:
		return controllers.GetWishesPage(ctx, request, db)
	case WishesPDFEndpoint:
		return controllers.GetWishesPDFPage(ctx, request, db)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, fmt.Errorf("endpoint %q not found", endpoint)
	}
}

func main() {
	lambda.Start(handler)
}

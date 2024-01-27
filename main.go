package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"wishy/controllers"
	"wishy/mongodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Connect to MongoDB
	client := mongodb.Connect(os.Getenv("MONGODB_URI"))
	db := client.Database(os.Getenv("DB_NAME"))

	endpoint := request.Path
	// Execute different logic based on the endpoint
	switch endpoint {
	case "/categories":
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
		return *response, nil
	case "/wishes":
		response, err := controllers.GetWishes(ctx, request, db)
		if err != nil {
			errInt, e := strconv.Atoi(err.Error())
			if e != nil {
				return events.APIGatewayProxyResponse{}, fmt.Errorf("error: %v - %v", err, e)
			}
			return events.APIGatewayProxyResponse{
				StatusCode: errInt,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, fmt.Errorf("get wishes: %v", err)
		}
		return *response, nil
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

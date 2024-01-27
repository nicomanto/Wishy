package common

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func JSONResponse(m interface{}) (*events.APIGatewayProxyResponse, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(j),
	}, nil
}

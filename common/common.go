package common

import (
	"encoding/json"
	"fmt"

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

func HTMLResponse(body string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       body,
	}, nil
}

func PDFResponse(body, filename string, forceDownload bool) (*events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Content-Type": "application/pdf",
		"Access-Control-Allow-Origin":  "*", // ✅ Allows all origins
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type"}
	if forceDownload {
		headers["Content-Disposition"] = "attachment; filename=" + filename
	} else {
		headers["Content-Disposition"] = "inline; filename=" + filename
	}
	fmt.Println("Encoded PDF length:", len(body)) // Should be > 1KB
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            body, // Only first 100 chars
		IsBase64Encoded: true, // Disable encoding
	}, nil
}

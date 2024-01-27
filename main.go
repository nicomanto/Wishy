package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (string, error) {
	return fmt.Sprintf("Hello from Golang Lambda!"), nil
}

func main() {
	lambda.Start(handler)
}

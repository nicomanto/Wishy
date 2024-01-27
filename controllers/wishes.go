package controllers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"wishy/common"
	"wishy/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWishes(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*events.APIGatewayProxyResponse, error) {
	wishes := []models.Wish{}
	cur, err := db.Collection(models.Wish{}.DBCollectionName()).Find(ctx, bson.M{})
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	if err := cur.All(ctx, &wishes); err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}

	// Create an HTML template
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Item List</title>
		</head>
		<body>
			<h1>Item List</h1>
			<ul>
				{{range .}}
					<li>
						<strong>ID:</strong> {{.ID}}<br>
						<strong>Name:</strong> {{.Name}}<br>
						<strong>Link:</strong> <a href="{{.Link}}" target="_blank">{{.Link}}</a><br>
						<strong>Category:</strong> {{.Cat}}<br>
					</li>
				{{end}}
			</ul>
		</body>
		</html>
	`)
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}

	var responseBody strings.Builder
	err = tmpl.Execute(&responseBody, wishes)
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	return common.HTMLResponse(responseBody.String())
}

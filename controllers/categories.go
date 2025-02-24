package controllers

import (
	"context"
	"fmt"
	"net/http"
	"wishy/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCategories(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) ([]models.Category, error) {
	categories := []models.Category{}
	cur, err := db.Collection(models.Category{}.DBCollectionName()).Find(ctx, bson.M{})
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	if err := cur.All(ctx, &categories); err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	return categories, nil
}

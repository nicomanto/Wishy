package controllers

import (
	"context"
	"wishy/common"
	"wishy/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCategeories(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*events.APIGatewayProxyResponse, error) {
	categories := []models.Category{}
	cur, err := db.Collection(models.Category{}.DBCollectionName()).Find(ctx, bson.M{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	if err := cur.All(ctx, &categories); err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	return common.JSONResponse(categories)
}

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

func GetWishes(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*events.APIGatewayProxyResponse, error) {
	wishes := []models.Wish{}
	cur, err := db.Collection(models.Wish{}.DBCollectionName()).Find(ctx, bson.M{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	if err := cur.All(ctx, &wishes); err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	return common.JSONResponse(wishes)
}

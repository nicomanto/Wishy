package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"wishy/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWishes(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*models.UserWishes, error) {
	// get uid
	uid, exists := request.QueryStringParameters["uid"]
	if !exists || uid == "" {
		logrus.Errorln(errors.New("require uid"))
		return nil, fmt.Errorf("%d", http.StatusBadRequest)
	}
	uidObjectId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		logrus.Errorln(fmt.Errorf("uid not valid: %v", err))
		return nil, fmt.Errorf("%d", http.StatusBadRequest)
	}
	// fetch the user
	user := models.User{}
	if err := db.Collection(user.DBCollectionName()).FindOne(ctx, bson.M{"_id": uidObjectId}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logrus.Errorln(fmt.Errorf("user %s: %v", uidObjectId.Hex(), err))
			return nil, fmt.Errorf("%d", http.StatusNotFound)
		}
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	// fecth wishes
	wishes := []models.WishByCategory{}
	cur, err := db.Collection(models.Wish{}.DBCollectionName()).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"uid": uidObjectId, "active": true}},
		{"$sort": bson.M{"preference": -1, "name": 1}},
		{"$group": bson.M{
			"_id": "$cat.name",
			"wishes": bson.M{
				"$push": bson.M{"name": "$name", "link": "$link", "preference": "$preference"},
			},
		}},
		{"$sort": bson.M{"_id": 1}},
	})
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	if err := cur.All(ctx, &wishes); err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	return &models.UserWishes{Wishes: wishes, Username: user.Name}, nil
}

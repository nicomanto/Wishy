package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"wishy/common"
	"wishy/models"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWishes(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*events.APIGatewayProxyResponse, error) {
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
	wishes := []models.UserWishes{}
	cur, err := db.Collection(models.Wish{}.DBCollectionName()).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"uid": uidObjectId}},
		{"$lookup": bson.M{
			"from":         "users",     // The users collection
			"localField":   "uid",       // Field in the wishes collection (uid)
			"foreignField": "_id",       // Field in the users collection (_id)
			"as":           "user_info", // Field to store the result of the lookup
		}},
		// Unwind the user_info array if you expect only one user per wish
		{"$unwind": "$user_info"},
		{"$sort": bson.M{"name": 1}},
		{"$group": bson.M{
			"_id": "$cat.name",
			"wishes": bson.M{
				"$push": bson.M{"name": "$name", "link": "$link"},
			},
		}},
		{"$addFields": bson.M{"username": "$user_info.name"}},
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

	// load html page
	var responseBody strings.Builder
	err = templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&responseBody, wishes[0])
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	return common.HTMLResponse(responseBody.String())
}

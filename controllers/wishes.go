package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"wishy/common"
	"wishy/models"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWishes(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*events.APIGatewayProxyResponse, error) {
	wishes := []struct {
		Cat    string `json:"cat" bson:"_id"`
		Wishes []struct {
			Name string `json:"name" bson:"name"`
			Link string `json:"link" bson:"link"`
		} `json:"wishes" bson:"wishes"`
	}{}
	cur, err := db.Collection(models.Wish{}.DBCollectionName()).Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id": "$cat.name",
			"wishes": bson.M{
				"$push": bson.M{"name": "$name", "link": "$link"},
			},
		}},
		{"$sort": bson.M{"_id": -1}},
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
	err = templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&responseBody, wishes)
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	return common.HTMLResponse(responseBody.String())
}

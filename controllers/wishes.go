package controllers

import (
	"context"
	"encoding/base64"
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

// GetWishesPage renders the user's public wishlist as HTML.
func GetWishesPage(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	templates.InitHtmlTpls()
	wishes, err := getUserWishes(ctx, request, db)
	if err != nil {
		return renderWishError(err)
	}

	var body strings.Builder
	if err := templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&body, wishes); err != nil {
		logrus.Errorln(err)
		return wishRenderErrorResponse(err)
	}
	response, err := common.HTMLResponse(body.String())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return *response, nil
}

// GetWishesPDFPage generates and returns the user's wishlist as a downloadable PDF.
func GetWishesPDFPage(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (events.APIGatewayProxyResponse, error) {
	templates.InitHtmlTpls()
	wishes, err := getUserWishes(ctx, request, db)
	if err != nil {
		return renderWishError(err)
	}

	pdfBytes, err := wishes.GenerateWishlistPDF()
	if err != nil {
		return renderWishError(err)
	}
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
	response, err := common.PDFResponse(pdfBase64, wishes.Username+"-wishlist.pdf", true)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return *response, nil
}

func renderWishError(cause error) (events.APIGatewayProxyResponse, error) {
	friendlyError := models.FriendlyErrorInit(cause.Error())
	var body strings.Builder
	if err := templates.HtmlTpls[templates.ErrorPageHtmlTemplateType].Execute(&body, friendlyError); err != nil {
		logrus.Errorln(err)
		return wishRenderErrorResponse(err)
	}
	response, err := common.HTMLResponse(body.String())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return *response, nil
}

func wishRenderErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, fmt.Errorf("render wishes response: %w", err)
}

// getUserWishes loads active wishes for the requested user and groups them by category.
func getUserWishes(ctx context.Context, request events.APIGatewayProxyRequest, db *mongo.Database) (*models.UserWishes, error) {
	// Validate the requested user ID before accessing the database.
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
	// Load the user to get display settings and the last wishlist update.
	user := models.User{}
	if err := db.Collection(user.DBCollectionName()).FindOne(ctx, bson.M{"_id": uidObjectId}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logrus.Errorln(fmt.Errorf("user %s: %v", uidObjectId.Hex(), err))
			return nil, fmt.Errorf("%d", http.StatusNotFound)
		}
		logrus.Errorln(err)
		return nil, fmt.Errorf("%d", http.StatusInternalServerError)
	}
	// Fetch active wishes and group them by category.
	wishes := []models.WishByCategory{}
	cur, err := db.Collection(models.Wish{}.DBCollectionName()).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"uid": uidObjectId, "active": true}},
		{"$sort": bson.D{{Key: "preference", Value: -1}, {Key: "name", Value: 1}}},
		{"$group": bson.M{
			"_id": "$cat.name",
			"wishes": bson.M{
				"$push": bson.M{"name": "$name", "link": "$link", "preference": "$preference", "ts": "$ts"},
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
	// Calculate the recent flag using the user's configured time window.
	for i := range wishes {
		for j := range wishes[i].Wishes {
			wishes[i].Wishes[j].SetRecent(user.IsRecentMonthConfig)
		}
	}
	return &models.UserWishes{Wishes: wishes, Username: user.Name, LastUpdate: user.LastWishUpdate.Format("02/01/2006")}, nil
}

package main

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
	"wishy/controllers"
	"wishy/models"
	"wishy/templates"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

var mockWishesList models.UserWishes
var mockFriendlyError models.FriendlyError

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

// tests setup
func setup() {
	// pre load html page
	templates.InitHtmlTpls()
	mockWishesList = models.UserWishes{
		Username: "username",
		Wishes: []models.WishByCategory{
			{
				Cat: "cat-1",
				Wishes: []models.BaseWish{
					{
						Name:       "wish-cat-1-name-1",
						Link:       "wish-cat-1-link-1",
						Preference: models.Low,
					},
					{
						Name:       "wish-cat-1-name-2",
						Link:       "wish-cat-1-link-2",
						Preference: models.Medium,
					},
					{
						Name:       "wish-cat-1-name-3",
						Link:       "wish-cat-1-link-3",
						Preference: models.High,
					},
				},
			},
			{
				Cat: "cat-2",
				Wishes: []models.BaseWish{
					{
						Name:       "wish-cat-2-name-1",
						Link:       "wish-cat-2-link-1",
						Preference: models.High,
					},
					{
						Name:       "wish-cat-2-name-2",
						Link:       "wish-cat-2-link-2",
						Preference: models.High,
					},
					{
						Name:       "wish-cat-2-name-3",
						Link:       "wish-cat-2-link-3",
						Preference: models.High,
					},
				}},
		},
	}
	mockFriendlyError = models.FriendlyErrorInit("400")
}

func TestRenderHtmlWishListPage(t *testing.T) {
	r := require.New(t)
	// load html page
	var responseBody strings.Builder
	err := templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&responseBody, mockWishesList)
	r.NoError(err)
}

func TestRenderHtmlErrorPage(t *testing.T) {
	r := require.New(t)
	// load html page
	var responseBody strings.Builder
	err := templates.HtmlTpls[templates.ErrorPageHtmlTemplateType].Execute(&responseBody, mockFriendlyError)
	r.NoError(err)
}

func TestGetUserWishesErrors(t *testing.T) {
	r := require.New(t)
	resp, err := controllers.GetUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{}}, nil)
	r.Nil(resp)
	r.Equal("400", err.Error())
	resp, err = controllers.GetUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"uid": ""}}, nil)
	r.Nil(resp)
	r.Equal("400", err.Error())
	resp, err = controllers.GetUserWishes(context.Background(), events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"uid": "aasdf"}}, nil)
	r.Nil(resp)
	r.Equal("400", err.Error())
}

func TestGenerateWishlistPDF(t *testing.T) {
	r := require.New(t)
	b, err := mockWishesList.GenerateWishlistPDF()
	r.NoError(err)
	r.NotNil(b)
}

func TestSetRecent(t *testing.T) {
	r := require.New(t)
	wish := models.Wish{
		Ts: time.Now().AddDate(0, -2, 0),
	}
	r.False(wish.IsRecent)
	wish.SetRecent(3)
	r.True(wish.IsRecent)
	wish.Ts = time.Now().AddDate(0, -4, 0)
	wish.SetRecent(3)
	r.False(wish.IsRecent)
}

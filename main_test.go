package main

import (
	"os"
	"strings"
	"testing"
	"wishy/templates"

	"github.com/stretchr/testify/require"
)

var mockWishesList []struct {
	Cat    string `json:"cat" bson:"_id"`
	Wishes []struct {
		Name string `json:"name" bson:"name"`
		Link string `json:"link" bson:"link"`
	} `json:"wishes" bson:"wishes"`
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

// tests setup
func setup() {
	// pre load html page
	templates.InitHtmlTpls()
	mockWishesList = []struct {
		Cat    string "json:\"cat\" bson:\"_id\""
		Wishes []struct {
			Name string "json:\"name\" bson:\"name\""
			Link string "json:\"link\" bson:\"link\""
		} "json:\"wishes\" bson:\"wishes\""
	}{
		{
			Cat: "cat-1",
			Wishes: []struct {
				Name string "json:\"name\" bson:\"name\""
				Link string "json:\"link\" bson:\"link\""
			}{
				{
					Name: "wish-cat-1-name-1",
					Link: "wish-cat-1-link-1",
				},
				{
					Name: "wish-cat-1-name-2",
					Link: "wish-cat-1-link-2",
				},
				{
					Name: "wish-cat-1-name-3",
					Link: "wish-cat-1-link-3",
				},
			},
		},
		{
			Cat: "cat-2",
			Wishes: []struct {
				Name string "json:\"name\" bson:\"name\""
				Link string "json:\"link\" bson:\"link\""
			}{
				{
					Name: "wish-cat-2-name-1",
					Link: "wish-cat-2-link-1",
				},
				{
					Name: "wish-cat-2-name-2",
					Link: "wish-cat-2-link-2",
				},
				{
					Name: "wish-cat-2-name-3",
					Link: "wish-cat-2-link-3",
				},
			},
		},
	}
}

func TestRenderWishListHtmlPage(t *testing.T) {
	r := require.New(t)
	// load html page
	var responseBody strings.Builder
	err := templates.HtmlTpls[templates.WishListHtmlTemplateType].Execute(&responseBody, mockWishesList)
	r.NoError(err)
}

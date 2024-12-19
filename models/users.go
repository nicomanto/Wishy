package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	LastWishUpdate time.Time          `json:"last_wish_update" bson:"last_wish_update"`
}

func (c User) DBCollectionName() string {
	return "users"
}

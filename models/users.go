package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	Name                string             `json:"name" bson:"name"`
	LastWishUpdate      time.Time          `json:"last_wish_update" bson:"last_wish_update"`
	IsRecentMonthConfig int                `json:"is_recent_month_config" bson:"is_recent_month_config"`
}

func (c User) DBCollectionName() string {
	return "users"
}

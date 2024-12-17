package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wish struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Link       string             `json:"link" bson:"link"`
	Category   BaseCategory       `json:"cat" bson:"cat"`
	UserId     primitive.ObjectID `json:"uid" bson:"uid"`
	Active     bool               `json:"active" bson:"active"`
	Preference PreferenceType     `json:"preference" bson:"preference"`
	Ts         time.Time          `json:"ts" bson:"ts"`
}

func (c Wish) DBCollectionName() string {
	return "wishes"
}

type PreferenceType int

const (
	Low PreferenceType = iota + 1
	Medium
	High
)

type WishByCategory struct {
	Cat    string `json:"cat" bson:"_id"`
	Wishes []struct {
		Name       string         `json:"name" bson:"name"`
		Link       string         `json:"link" bson:"link"`
		Preference PreferenceType `json:"preference" bson:"preference"`
		Ts         time.Time      `json:"ts" bson:"ts"`
	} `json:"wishes" bson:"wishes"`
}

type UserWishes struct {
	Wishes     []WishByCategory `json:"wishes" bson:"wishes"`
	Username   string           `json:"username" bson:"username"`
	LastUpdate string           `json:"last_update" bson:"last_update"`
}

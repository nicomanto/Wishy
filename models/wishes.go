package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wish struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Link     string             `json:"link" bson:"link"`
	Category BaseCategory       `json:"cat" bson:"cat"`
	UserId   primitive.ObjectID `json:"uid" bson:"uid"`
}

func (c Wish) DBCollectionName() string {
	return "wishes"
}

type WishByCategory struct {
	Cat    string `json:"cat" bson:"_id"`
	Wishes []struct {
		Name string `json:"name" bson:"name"`
		Link string `json:"link" bson:"link"`
	} `json:"wishes" bson:"wishes"`
}

type UserWishes struct {
	Wishes   []WishByCategory `json:"wishes" bson:"wishes"`
	Username string           `json:"username" bson:"username"`
}

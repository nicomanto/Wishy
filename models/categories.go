package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	BaseCategory `json:",inline" bson:",inline"`
	UserId       primitive.ObjectID `json:"uid" bson:"uid"`
}

type BaseCategory struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

func (c Category) DBCollectionName() string {
	return "categories"
}

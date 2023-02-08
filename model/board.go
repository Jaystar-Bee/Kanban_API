package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	UserID string             `json:"user_id" bson:"user_id"`
}

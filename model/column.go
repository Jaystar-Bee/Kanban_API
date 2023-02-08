package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Column struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Color   string             `json:"color" bson:"color"`
	BoardID primitive.ObjectID `json:"board_id" bson:"board_id"`
	UserID  string `json:"user_id" bson:"user_id"`
}

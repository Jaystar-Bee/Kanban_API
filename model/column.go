package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Column struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Name string `json:"name" bson:"name" validate:"required, min=2, max=20"`

	Color string `json:"color" bson:"color" validate:"required, min=2, max=10"`
}

type ColumnsReply struct {
	Body []*Column
}

type ColumnParam struct {
	Body []*Column
}

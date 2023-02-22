package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Column is the model for column
//
//	swagger:model
type Column struct {
	// ID for the column
	//
	// required: true
	//
	// example: 5e4d3b5b6b6b6b6b6b6b6b6b
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// Name of the column
	//
	// required: true
	//
	// example: My Column
	Name string `json:"name" bson:"name" validate:"required, min=2, max=20"`
	// Color of the column
	//
	// required: true
	//
	// example: #000000
	Color string `json:"color" bson:"color" validate:"required, min=2, max=10"`
}

// // Board ID of the column
// //
// // required: true
// //
// // example: 5e4d3b5b6b6b6b6b6b6b6b6b
// BoardID primitive.ObjectID `json:"board_id" bson:"board_id"`
// // User ID of the column
// //
// // required: true
// //
// // example: 5e4d3b5b6b6b6b6b6b6b6b6b
// UserID string `json:"user_id" bson:"user_id"`

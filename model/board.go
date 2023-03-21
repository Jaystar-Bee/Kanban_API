package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model
type Board struct {
	// The id for the board
	// required: true
	// example: 5f1f9b9e0f1c9c0001e1b1b1
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// The name for the board
	// required: true
	// example: My Board
	Name string `json:"name" bson:"name" validate:"required, min=2, max=20"`
	// Columns for the board
	//
	// Extensions:
	// ---
	// x-property-value: value
	// x-property-array:
	// - value1
	// - value2
	// x-property-object:
	//   key1: value1
	//   key2: value2
	// ---
	Columns []Column `json:"columns" bson:"columns"`
	// The user id for the board
	// required: true
	// example: 5f1f9b9e0f1c9c0001e1b1b1
	UserID string `json:"user_id" bson:"user_id"`
}

// swagger:model BoardRequest
type BoardRequest struct {
	// The name for the board
	// required: true
	// example: My Board
	Name string `json:"name" bson:"name" validate:"required, min=2, max=20"`
	// Columns for the board
	//
	// Extensions:
	// ---
	// x-property-value: value
	// x-property-array:
	// - value1
	// - value2
	// x-property-object:
	//   key1: value1
	//   key2: value2
	// ---
	Columns []Column `json:"columns" bson:"columns"`
}

// swagger:response
type BoardReply struct {
	// in: body
	Body *Board
}

// swagger:parameters createBoard
type BoardParam struct {
	// in: body
	Body *BoardRequest
}

// swagger:response
type ErrorResponse struct {
	// The error message
	// in: body
	Body struct {
		// The error message
		// example: Error message
		Message string `json:"message"`
	}
}

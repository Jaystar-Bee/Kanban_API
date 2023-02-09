package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Board is the model for board
//
//	swagger:model
type Board struct {
	// ID for the user
	//
	// required: true
	//
	// example: 5e4d3b5b6b6b6b6b6b6b6b6b
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// Name of the board
	//
	// required: true
	//
	// example: My Board
	Name string `json:"name" bson:"name"`
	// Color of the board
	//
	// required: true
	//
	Columns []Column `json:"columns" bson:"columns"`
	// User ID of the board
	//
	// required: true
	//
	// example: 5e4d3b5b6b6b6b6b6b6b6b6b
	UserID string `json:"user_id" bson:"user_id"`
}

// Board reply is the model for board reply
//
//	swagger:response
type BoardReply struct {
	// in: body
	Body *Board
}

//Error reply is the model for error reply
//
//	swagger:response
type ErrorReply struct {
	// in: body
	Body *error
}

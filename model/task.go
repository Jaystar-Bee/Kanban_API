package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title" validate:"required, min=2"`
	Description string             `json:"description" bson:"description" validate:"required, min=2"`
	Status      string             `json:"status" bson:"status" validate:"required, min=2, max=20"`
	SubTasks    []SubTask          `json:"sub_tasks" bson:"sub_tasks" validate:"required"`
	BoardID     primitive.ObjectID `json:"board_id" bson:"board_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
}

type SubTask struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" `
	Title  string             `json:"title" bson:"title" validate:"required, min=2, max=20"`
	IsDone bool               `json:"is_done" bson:"is_done"`
}

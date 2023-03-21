package handlers

import (
	"context"
	helpers "kanban-task/helper"
	"kanban-task/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ColumnHandler struct {
	ctx             context.Context
	boardCollection *mongo.Collection
	redisClient     *redis.Client
}

func NewColumnHandler(ctx context.Context, boardCollection *mongo.Collection, redisClient *redis.Client) *ColumnHandler {
	return &ColumnHandler{
		ctx:             ctx,
		boardCollection: boardCollection,
		redisClient:     redisClient,
	}
}

// List all Columns

// swagger:route GET /boards/{id}/columns Columns ListColumn
//
// This route get all the columns for a board
//
// Get all the columns for a board
//
// Produces:
// -application/json
//
// Parameters:
// + name: Authorization
//   in: header
//   description: "Authorization token"
//   required: true
//   type: string
// + name: id
//   in: path
//   description: "Board ID"
//   required: true
//   type: string
//
// Responses:
// 200: ColumnReply
// 500: ErrorResponse

func (handler *ColumnHandler) ListColumnHandler(c *gin.Context) {
	var board model.Board
	boardID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = handler.boardCollection.FindOne(handler.ctx, bson.D{{"_id", boardID}}).Decode(&board)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, board.Columns)
}

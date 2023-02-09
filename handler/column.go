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
	taskCollection  *mongo.Collection
	redisClient     *redis.Client
}

func NewColumnHandler(ctx context.Context, boardCollection *mongo.Collection, taskCollection *mongo.Collection, redisClient *redis.Client) *ColumnHandler {
	return &ColumnHandler{
		ctx:             ctx,
		boardCollection: boardCollection,
		taskCollection:  taskCollection,
		redisClient:     redisClient,
	}
}

// List all Columns
func (handler *ColumnHandler) ListColumnHandler(c *gin.Context) {
	var board model.Board
	boardID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	err = handler.boardCollection.FindOne(handler.ctx, bson.D{{"_id", boardID}}).Decode(&board)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	// var columns []model.Column
	// columns = board.Columns
	// fmt.Println(columns)
	c.JSON(http.StatusOK, board.Columns)
}

// Create Column

// func (handler *ColumnHandler) InsertColumnHandler(c *gin.Context) {
// 	var column model.Column
// 	boardID, err := primitive.ObjectIDFromHex(c.Param("id"))
// 	err = c.ShouldBind(&column)

// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	result, err := handler.columnCollection.InsertOne(handler.ctx, column)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	column.ID = primitive.NewObjectID()
// 	column.BoardID = boardID
// 	column.UserID = "UI5f9f1b9c1c9d440000a1e1f1"
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": result,
// 		"column":  column,
// 	})
// }

// Delete Column

// func (handler *ColumnHandler) DeleteColumnHandler(c *gin.Context) {
// 	columnID, err := helpers.ToPrimitive(c.Param("id"))
// 	// columnID, err := primitive.ObjectIDFromHex(c.Param("id"))
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	res, err := handler.columnCollection.DeleteOne(handler.ctx, bson.M{"_id": columnID})
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// GET a Column

// func (handler *ColumnHandler) GetColumnHandler(c *gin.Context) {
// 	columnID, err := helpers.ToPrimitive(c.Param("id"))
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	var column model.Column
// 	err = handler.columnCollection.FindOne(handler.ctx, bson.M{"_id": columnID}).Decode(&column)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, column)

// }

// Edit a column

// func (handler *ColumnHandler) UpdateColumnHandler(c *gin.Context) {
// 	columnID, err := helpers.ToPrimitive(c.Param("id"))
// 	var column model.Column
// 	err = c.ShouldBindJSON(&column)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	element := bson.D{{"_id", columnID}}
// 	filter := bson.D{{"$set", bson.D{
// 		{"name", column.Name},
// 		{"color", column.Color},
// 	}}}

// 	update, err := handler.columnCollection.UpdateOne(handler.ctx, element, filter)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"message": err.Error(),
// 		})
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"result": update,
// 		"column": column,
// 	})
// }

package handlers

import (
	"context"
	"net/http"

	"kanban-task/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var err error

type BoardHandler struct {
	ctx              context.Context
	boardCollection  *mongo.Collection
	columnCollection *mongo.Collection
	taskCollection   *mongo.Collection
	redisClient      *redis.Client
}

func NewBoardHandler(
	ctx context.Context,
	boardCollection *mongo.Collection,
	columnCollection *mongo.Collection,
	taskCollection *mongo.Collection,
	redisClient *redis.Client,
) *BoardHandler {
	return &BoardHandler{
		ctx:              ctx,
		boardCollection:  boardCollection,
		columnCollection: boardCollection,
		taskCollection:   boardCollection,
		redisClient:      redisClient,
	}
}

//swagger:route GET /boards Boards GetAllBoards
//
// This route get all the boards created by the user
//
// Produces
// -application/json
//
// Responses:
//	200: model.Board
//	500:
// "message": err
//

func (handler *BoardHandler) ListBoardHandler(c *gin.Context) {
	var boards []model.Board
	cursor, err := handler.boardCollection.Find(handler.ctx, bson.D{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = cursor.All(handler.ctx, &boards); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, boards)
}

func (handler *BoardHandler) InsertBoardHandler(c *gin.Context) {
	var board model.Board
	if err := c.ShouldBindJSON(&board); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	board.UserID = "UI5f9f1b9c1c9d440000a1e1f1"
	board.ID = primitive.NewObjectID()
	_, err := handler.boardCollection.InsertOne(handler.ctx, board)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board)
}

func (handler *BoardHandler) GetBoard(c *gin.Context) {

	var id = c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	var board model.Board
	err := handler.boardCollection.FindOne(handler.ctx, bson.M{"_id": objectID}).Decode(&board)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board)
}
func (handler *BoardHandler) DeleteBoard(c *gin.Context) {
	var id = c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)

	board, err := handler.boardCollection.DeleteOne(handler.ctx, bson.M{"_id": objectID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	column, err := handler.columnCollection.DeleteMany(handler.ctx, bson.M{"board_id": objectID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	task, err := handler.taskCollection.DeleteMany(handler.ctx, bson.M{"board_id": objectID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Board deleted successfully",
		"boardMessage":  board,
		"columnMessage": column,
		"taskMessage":   task,
	})

}

func (handler *BoardHandler) UpdateBoard(c *gin.Context) {
	id := c.Param("id")
	var board model.Board
	objectID, _ := primitive.ObjectIDFromHex(id)
	err := c.ShouldBind(&board)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	update := bson.D{{"$set", bson.D{{"name", board.Name}}}}

	res, err := handler.boardCollection.UpdateOne(handler.ctx, bson.D{{"_id", objectID}}, update)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})

}

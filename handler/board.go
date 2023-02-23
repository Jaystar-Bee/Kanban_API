package handlers

import (
	"context"
	"net/http"

	helpers "kanban-task/helper"
	"kanban-task/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var err error

type BoardHandler struct {
	ctx             context.Context
	boardCollection *mongo.Collection
	taskCollection  *mongo.Collection
	redisClient     *redis.Client
}

func NewBoardHandler(
	ctx context.Context,
	boardCollection *mongo.Collection,
	taskCollection *mongo.Collection,
	redisClient *redis.Client,
) *BoardHandler {
	return &BoardHandler{
		ctx:             ctx,
		boardCollection: boardCollection,
		taskCollection:  taskCollection,
		redisClient:     redisClient,
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
//	200: BoardReply
//	500: ErrorReply
//

func (handler *BoardHandler) ListBoardHandler(c *gin.Context) {
	var boards []model.Board
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	var userID = user.(*helpers.Claims).UserID
	cursor, err := handler.boardCollection.Find(handler.ctx, bson.M{"user_id": userID})
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
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	board.UserID = user.(*helpers.Claims).UserID
	board.ID = primitive.NewObjectID()
	var columns []model.Column

	for _, column := range board.Columns {
		column.ID = primitive.NewObjectID()
		columns = append(columns, column)
	}
	board.Columns = columns

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

	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.(*helpers.Claims).UserID

	// getting board from the database
	err := handler.boardCollection.FindOne(handler.ctx, bson.M{"_id": objectID, "user_id": userID}).Decode(&board)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board)
}

func (handler *BoardHandler) DeleteBoard(c *gin.Context) {
	var id = c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)

	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.(*helpers.Claims).UserID

	// deleting board from the database
	board, err := handler.boardCollection.DeleteOne(handler.ctx, bson.M{"_id": objectID, "user_id": userID})
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
		"message":      "Board deleted successfully",
		"boardMessage": board,
		"taskMessage":  task,
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
		return
	}

	var columns []model.Column

	for _, column := range board.Columns {
		if column.ID == primitive.NilObjectID {
			column.ID = primitive.NewObjectID()
		}
		columns = append(columns, column)
	}
	board.Columns = columns

	update := bson.D{{"$set", bson.D{
		{"name", board.Name},
		{"Columns", board.Columns},
	}}}

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

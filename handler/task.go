package handlers

import (
	"context"
	helpers "kanban-task/helper"
	"kanban-task/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskHandler struct {
	ctx         context.Context
	collection  *mongo.Collection
	redisClient *redis.Client
}

func NewTaskHandler(ctx context.Context,
	collection *mongo.Collection,
	redisClient *redis.Client) *TaskHandler {
	return &TaskHandler{
		ctx:         ctx,
		collection:  collection,
		redisClient: redisClient,
	}
}

//Get All Tasks

func (handler *TaskHandler) ListTaskHandler(c *gin.Context) {
	var tasks []model.Task
	cursor, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	err = cursor.All(handler.ctx, &tasks)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, tasks)

}

// Create a Task

func (handler *TaskHandler) InsertTaskHandler(c *gin.Context) {
	boardID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	var task model.Task
	err = c.ShouldBindJSON(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	task.BoardID = boardID
	task.UserID = c.Param("id")
	task.ID = primitive.NewObjectID()
	result, err := handler.collection.InsertOne(handler.ctx, task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
		"task":    task,
	})
}

// Get A Task

func (handler *TaskHandler) GetTaskHandler(c *gin.Context) {
	var task model.Task
	taskID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	err = handler.collection.FindOne(handler.ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, task)
}

//Delete Task

func (handler *TaskHandler) DeleteTaskHandler(c *gin.Context) {
	taskID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	result, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": taskID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, result)
}

// Update a Task

func (handler *TaskHandler) UpdateTaskHandler(c *gin.Context) {
	var task model.Task
	taskID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	err = c.ShouldBindJSON(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	filter := bson.M{"_id": taskID}

	options := bson.M{"$set": bson.D{
		{"title", task.Title},
		{"description", task.Description},
		{"status", task.Status},
		{"sub_tasks", task.SubTasks},
	}}

	result, err := handler.collection.UpdateOne(handler.ctx, filter, options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
		"task":    task,
	})
}

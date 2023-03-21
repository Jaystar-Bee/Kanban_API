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

// swagger:route GET /tasks Tasks GetAllTasks
//
// This route get all the tasks created by the user
//
// Get all the tasks for the user
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
//
// Responses:
// 200: TaskReply
// 500: ErrorResponse

func (handler *TaskHandler) ListTaskHandler(c *gin.Context) {
	var tasks []model.Task
	// Get board ID
	boardID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.(*helpers.Claims).UserID

	// Get all tasks
	cursor, err := handler.collection.Find(handler.ctx, bson.M{"user_id": userID, "board_id": boardID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = cursor.All(handler.ctx, &tasks)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tasks)

}

// Create a Task

// swagger:route POST /tasks/:id Tasks CreateTask
//
// This route create a task
//
// Create a task
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
// + name: task
//   in: body
//   description: "Task to be created"
//   required: true
//   type: TaskRequest
//   schema:
//     "$ref": "#/definitions/TaskRequest"
//
// Responses:
// 200: TaskReply
// 500: ErrorResponse

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
	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	task.UserID = user.(*helpers.Claims).UserID
	task.ID = primitive.NewObjectID()
	// Give subTask ID
	var subTasks []model.SubTask
	for _, subTask := range task.SubTasks {
		if subTask.ID == primitive.NilObjectID {
			subTask.ID = primitive.NewObjectID()
		}
		subTasks = append(subTasks, subTask)
	}
	task.SubTasks = subTasks

	// Insert Task
	result, err := handler.collection.InsertOne(handler.ctx, task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
		"task":    task,
	})
}

// Get A Task

// swagger:route GET /tasks/:id Tasks GetTask
//
// This route get a task
//
// Get a task
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
//   description: "Task ID"
//   required: true
//   type: string
//
// Responses:
// 200: TaskReply
// 500: ErrorResponse

func (handler *TaskHandler) GetTaskHandler(c *gin.Context) {
	var task model.Task
	taskID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.(*helpers.Claims).UserID

	// Get Task
	err = handler.collection.FindOne(handler.ctx, bson.M{"_id": taskID, "user_id": userID}).Decode(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, task)
}

//Delete Task

// swagger:route DELETE /tasks/:id Tasks DeleteTask
//
// This route delete a task
//
// Delete a task
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
//   description: "Task ID"
//   required: true
//   type: string
//
// Responses:
// 200: TaskReply
// 500: ErrorResponse

func (handler *TaskHandler) DeleteTaskHandler(c *gin.Context) {
	taskID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.(*helpers.Claims).UserID

	// Delete Task
	result, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": taskID, "user_id": userID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Update a Task

// swagger:route PUT /tasks/:id Tasks UpdateTask
//
// This route update a task
//
// Update a task
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
//   description: "Task ID"
//   required: true
//   type: string
// + name: task
//   in: body
//   description: "Task"
//   required: true
//   type: TaskRequest
//   schema:
//     "$ref": "#/definitions/TaskRequest"
//
// Responses:
// 200: TaskReply
// 500: ErrorResponse

func (handler *TaskHandler) UpdateTaskHandler(c *gin.Context) {
	var task model.Task
	taskID, err := helpers.ToPrimitive(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = c.ShouldBindJSON(&task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	var subTasks []model.SubTask
	for _, subTask := range task.SubTasks {
		if subTask.ID == primitive.NilObjectID {
			subTask.ID = primitive.NewObjectID()
		}
		subTasks = append(subTasks, subTask)
	}
	task.SubTasks = subTasks

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
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": result,
		"task":    task,
	})
}

//Delete Task By Status

// swagger:route DELETE /tasks/:status Tasks DeleteTask
//
// This route delete task with a particular status
//
// Delete the task with the specified status
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
// + name: status
//   in: path
//   description: "Task status"
//   required: true
//   type: string
// + name: id
//   in: path
//   description: "User Id"
//   required: true
//   type: string
//
// Responses:
// 200: TaskReply
// 500: ErrorResponse

func (handler *TaskHandler) DeleteTasksByStatus(c *gin.Context) {
	status := c.Param("status")

	// getting user ID from the header
	user, exist := c.Get("user")
	if !exist {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.(*helpers.Claims).UserID

	result, err := handler.collection.DeleteMany(handler.ctx, bson.M{"status": status, "user_id": userID})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)

}

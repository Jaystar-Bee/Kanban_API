//	KanBan API
//
//	Kanban API is an restful API for task management for developers or any corperate workers
//
// It is a standalone API
//
//	Schemes: http, https
//	Host: localhost
//	BasePath: /api/v1/
//	Version: 1.0.0
//	Contact: John Ayilara <jbayilara@gmail.com>	https://bolu.netlify.app
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	swagger:meta
package main

import (
	"context"
	"fmt"

	handlers "kanban-task/handler"
	initializer "kanban-task/initial"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	PORT             = ":3000"
	boardCollection  *mongo.Collection
	columnCollection *mongo.Collection
	taskCollection   *mongo.Collection
	ctx              context.Context
	redisClient      *redis.Client
	boardHandler     *handlers.BoardHandler
	columnHandler    *handlers.ColumnHandler
	taskHandler      *handlers.TaskHandler
)

func init() {
	ctx, boardCollection, columnCollection, taskCollection = initializer.MongoConnect()
	redisClient = initializer.RedisConnect()
	boardHandler = handlers.NewBoardHandler(ctx, boardCollection, columnCollection, taskCollection, redisClient)
	columnHandler = handlers.NewColumnHandler(ctx, columnCollection, taskCollection, redisClient)
	taskHandler = handlers.NewTaskHandler(ctx, taskCollection, redisClient)
}

func main() {
	router := gin.Default()
	initRoute := router.Group("/api/v1")

	boardRoute := initRoute.Group("/boards")
	boardRoute.GET("/", boardHandler.ListBoardHandler)
	boardRoute.POST("/", boardHandler.InsertBoardHandler)
	boardRoute.GET("/:id", boardHandler.GetBoard)
	boardRoute.DELETE("/:id", boardHandler.DeleteBoard)
	boardRoute.PUT("/:id", boardHandler.UpdateBoard)

	columnRoute := initRoute.Group("/columns")
	columnRoute.GET("/", columnHandler.ListColumnHandler)
	columnRoute.POST("/:id", columnHandler.InsertColumnHandler)
	columnRoute.GET("/:id", columnHandler.GetColumnHandler)
	columnRoute.DELETE("/:id", columnHandler.DeleteColumnHandler)
	columnRoute.PUT("/:id", columnHandler.UpdateColumnHandler)

	taskRoute := initRoute.Group("/tasks")
	taskRoute.GET("/", taskHandler.ListTaskHandler)
	taskRoute.POST("/:id", taskHandler.InsertTaskHandler)
	taskRoute.GET("/:id", taskHandler.GetTaskHandler)
	taskRoute.DELETE("/:id", taskHandler.DeleteTaskHandler)
	taskRoute.PUT("/:id", taskHandler.UpdateTaskHandler)

	router.Run(PORT)
	fmt.Println("Serve running on port: ", PORT)
}

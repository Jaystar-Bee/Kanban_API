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
	"log"

	"kanban-task/auths"
	handlers "kanban-task/handler"
	initializer "kanban-task/initial"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	PORT            = ":3000"
	boardCollection *mongo.Collection
	userCollection  *mongo.Collection
	taskCollection  *mongo.Collection
	ctx             context.Context
	redisClient     *redis.Client
	boardHandler    *handlers.BoardHandler
	columnHandler   *handlers.ColumnHandler
	authHandler     *handlers.AuthHandler
	taskHandler     *handlers.TaskHandler
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load .env")
	}

	ctx, boardCollection, userCollection, taskCollection = initializer.MongoConnect()
	redisClient = initializer.RedisConnect()
	boardHandler = handlers.NewBoardHandler(ctx, boardCollection, taskCollection, redisClient)
	columnHandler = handlers.NewColumnHandler(ctx, boardCollection, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, userCollection, redisClient)
	taskHandler = handlers.NewTaskHandler(ctx, taskCollection, redisClient)
}

func main() {
	router := gin.Default()
	initRoute := router.Group("/api/v1")
	userRoute := initRoute.Group("/users")
	userRoute.POST("/signup", authHandler.SignUp)
	userRoute.POST("/signin", authHandler.SignIn)

	initRoute.Use(auths.AuthMiddleware(redisClient))
	initRoute.POST("/users/logout", authHandler.Logout)
	boardRoute := initRoute.Group("/boards")
	boardRoute.GET("/", boardHandler.ListBoardHandler)
	boardRoute.POST("/", boardHandler.InsertBoardHandler)
	boardRoute.GET("/:id", boardHandler.GetBoard)
	boardRoute.DELETE("/:id", boardHandler.DeleteBoard)
	boardRoute.PUT("/:id", boardHandler.UpdateBoard)

	columnRoute := initRoute.Group("/columns")
	columnRoute.GET("/:id", columnHandler.ListColumnHandler) // board id

	taskRoute := initRoute.Group("/tasks")
	taskRoute.GET("/:id", taskHandler.ListTaskHandler) // board id
	taskRoute.POST("/:id", taskHandler.InsertTaskHandler)
	taskRoute.GET("/task/:id", taskHandler.GetTaskHandler)
	taskRoute.DELETE("/:id", taskHandler.DeleteTaskHandler)
	taskRoute.PUT("/:id", taskHandler.UpdateTaskHandler)

	router.Run(PORT)
	fmt.Println("Serve running on port: ", PORT)
}

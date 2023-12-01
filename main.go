//	KanBan API
//
//	Kanban API is an restful API for task management for developers or any corperate workers
//
// It is a standalone API
//
//	Schemes: http, https
//	Host: https://kanban-api-app-jaystar-bee-dev.apps.sandbox-m3.1530.p1.openshiftapps.com
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

// "log"
import (
	"context"
	"fmt"
	"log"

	"kanban-task/auths"
	handlers "kanban-task/handler"
	initializer "kanban-task/initial"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	PORT              = ":3000"
	boardCollection   *mongo.Collection
	userCollection    *mongo.Collection
	expiresCollection *mongo.Collection
	taskCollection    *mongo.Collection
	ctx               context.Context
	redisClient       *redis.Client
	boardHandler      *handlers.BoardHandler
	columnHandler     *handlers.ColumnHandler
	authHandler       *handlers.AuthHandler
	taskHandler       *handlers.TaskHandler
)

func init() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Unable to load .env file")
	// }

	ctx, boardCollection, userCollection, taskCollection, expiresCollection = initializer.MongoConnect()
	redisClient = initializer.RedisConnect()
	boardHandler = handlers.NewBoardHandler(ctx, boardCollection, taskCollection, redisClient)
	columnHandler = handlers.NewColumnHandler(ctx, boardCollection, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, userCollection, expiresCollection, redisClient)
	taskHandler = handlers.NewTaskHandler(ctx, taskCollection, redisClient)
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5173/"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
	}))

	initRoute := router.Group("/api/v1")
	userRoute := initRoute.Group("/users")
	userRoute.POST("/signup", authHandler.SignUp)
	userRoute.POST("/signin", authHandler.SignIn)

	initRoute.Use(auths.AuthMiddleware(ctx, redisClient, expiresCollection))
	initRoute.POST("/users/logout", authHandler.Logout)
	boardRoute := initRoute.Group("/boards")
	boardRoute.GET("", boardHandler.ListBoardHandler)
	boardRoute.POST("/", boardHandler.InsertBoardHandler)
	boardRoute.GET("/:id", boardHandler.GetBoard)
	boardRoute.DELETE("/:id", boardHandler.DeleteBoard)
	boardRoute.PUT("/:id", boardHandler.UpdateBoard)

	columnRoute := initRoute.Group("/columns")
	columnRoute.GET("/:id", columnHandler.ListColumnHandler) // board id
	columnRoute.DELETE("/:status", taskHandler.DeleteTasksByStatus)

	taskRoute := initRoute.Group("/tasks")
	taskRoute.GET("/:id", taskHandler.ListTaskHandler) // board id
	taskRoute.POST("/:id", taskHandler.InsertTaskHandler)
	taskRoute.GET("/task/:id", taskHandler.GetTaskHandler)
	taskRoute.DELETE("/:id", taskHandler.DeleteTaskHandler)
	taskRoute.PUT("/:id", taskHandler.UpdateTaskHandler)

	err := router.Run(PORT)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Serve running on port: ", PORT)
}

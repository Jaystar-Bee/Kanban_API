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

// "log"
import (
	"context"
	"fmt"
	"time"

	"kanban-task/auths"
	handlers "kanban-task/handler"
	initializer "kanban-task/initial"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

// "github.com/joho/godotenv"

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
	// 	log.Println("Unable to load .env file")
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

	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://127.0.0.1:5173", "http://localhost:5173", "http://kanbantask.netlify.app", "https://kanbantask.netlify.app"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials", "Access-Control-Max-Age"},

			MaxAge: 12 * time.Hour,
		}))

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
	columnRoute.DELETE("/:status", taskHandler.DeleteTasksByStatus)

	taskRoute := initRoute.Group("/tasks")
	taskRoute.GET("/:id", taskHandler.ListTaskHandler) // board id
	taskRoute.POST("/:id", taskHandler.InsertTaskHandler)
	taskRoute.GET("/task/:id", taskHandler.GetTaskHandler)
	taskRoute.DELETE("/:id", taskHandler.DeleteTaskHandler)
	taskRoute.PUT("/:id", taskHandler.UpdateTaskHandler)

	router.Run(PORT)
	fmt.Println("Serve running on port: ", PORT)
}

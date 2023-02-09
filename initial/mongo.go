package initializer

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoConnect() (ctx context.Context, boardCollection *mongo.Collection, columnCollection *mongo.Collection, taskCollection *mongo.Collection) {

	ctx = context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// connect mongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	boardCollection = client.Database("kanban").Collection("board")
	columnCollection = client.Database("kanban").Collection("column")
	taskCollection = client.Database("kanban").Collection("task")

	return ctx, boardCollection, columnCollection, taskCollection
}

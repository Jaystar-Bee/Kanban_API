package helpers

import (
	"context"
	"crypto/sha256"
	"kanban-task/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByName(collection *mongo.Collection, userName string, ctx context.Context) (model.User, error) {
	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": userName}).Decode(&user)
	return user, err
}

func CheckPassword(collection *mongo.Collection, userName string, password string, ctx context.Context) (model.User, error) {
	hashPassword := HashPassword(password)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": userName, "password": hashPassword}).Decode(&user)
	return user, err
}

func HashPassword(password string) string {
	h := sha256.New()
	newPassword := string(h.Sum([]byte(password)))
	return newPassword
}

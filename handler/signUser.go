package handlers

import (
	"context"
	"fmt"
	helpers "kanban-task/helper"
	"kanban-task/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	ctx         context.Context
	collection  *mongo.Collection
	redisClient *redis.Client
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		ctx:         ctx,
		collection:  collection,
		redisClient: redisClient,
	}

}

func (handler *AuthHandler) SignUp(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	_, err = helpers.GetUserByName(handler.collection, user.Username, handler.ctx)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User already Exist",
		})
		return
	}
	user.Password = helpers.HashPassword(user.Password)
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UserID = user.ID.Hex()
	_, err = handler.collection.InsertOne(handler.ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fullDetail, err := helpers.GetUserByName(handler.collection, user.Username, handler.ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "User created but you're unable to log in yet",
		})
		return
	}
	userDetails, err := helpers.GenerateToken(fullDetail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userDetails)

}

func (handler *AuthHandler) SignIn(c *gin.Context) {
	var user model.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	signedUser, err := helpers.CheckPassword(handler.collection, user.Username, user.Password, handler.ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	userResult, err := helpers.GenerateToken(signedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, userResult)
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"message": "The token is required",
		})
	}
	user, _ := c.Get("user")
	fmt.Println(user.(*helpers.Claims).Username)
	expireTime := user.(*helpers.Claims).StandardClaims.ExpiresAt

	reqToken := strings.Split(token, "Bearer ")[1]
	secs := time.Duration(expireTime) * time.Second
	err = handler.redisClient.Set(reqToken, "Invalid", secs).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unabale to logout at the moment",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out",
	})
}

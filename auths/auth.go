package auths

import (
	"context"
	"fmt"
	helpers "kanban-task/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type tokenStruct struct {
	Token string `bson:"token"`
}

func AuthMiddleware(ctx context.Context, redisClient *redis.Client, expires *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var realToken tokenStruct
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"message": "The token is required",
			})
		}
		reqToken := strings.Split(token, "Bearer ")[1]

		err := expires.FindOne(ctx, bson.M{"token": reqToken}).Decode(&tokenStruct{})
		fmt.Println(err.Error())
		if err.Error() != "mongo: no documents in result" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "An Invalid token",
			})
			c.Abort()
		}

		//
		// val, _ := redisClient.Get(reqToken).Result()
		// if val == "Invalid" {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 		"message": "An Invalid token",
		// 	})
		// 	c.Abort()
		// }
		//
		tokenString, claims, err := helpers.ValidateToken(reqToken)

		if err != nil || !tokenString.Valid || tokenString == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}

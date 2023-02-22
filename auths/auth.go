package auths

import (
	"fmt"
	helpers "kanban-task/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"message": "The token is required",
			})
		}
		reqToken := strings.Split(token, "Bearer ")[1]
		val, err := redisClient.Get("reqToken").Result()

		if val != "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "An Invalid token",
			})
		}

		tokenString, claims, err := helpers.ValidateToken(reqToken)

		if err != nil || !tokenString.Valid || tokenString == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set("user", claims)
		fmt.Println(claims)
		c.Next()
	}
}

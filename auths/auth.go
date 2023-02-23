package auths

import (
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
		val, _ := redisClient.Get(reqToken).Result()
		if val == "Invalid" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "An Invalid token",
			})
			c.Abort()
		}
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

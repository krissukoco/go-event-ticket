package auth

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(getterFunc func(string) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}
		split := strings.Split(authHeader, " ")
		if len(split) != 2 {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid authorization header"})
			return
		}
		if split[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid authorization method"})
			return
		}
		token := split[1]
		log.Println(token)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
			return
		}
		// Verify token
		userId, err := getterFunc(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}

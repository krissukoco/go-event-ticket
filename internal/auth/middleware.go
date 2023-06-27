package auth

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getUserIdFromToken(token string, secret string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	// Validate expiry
	if !t.Valid {
		return "", ErrInvalidToken
	}
	exp, err := t.Claims.GetExpirationTime()
	if err != nil {
		return "", ErrInvalidToken
	}
	if exp.Before(time.Now()) {
		return "", ErrExpiredToken
	}
	iss, err := t.Claims.GetIssuer()
	if err != nil {
		return "", ErrInvalidToken
	}
	if iss != issuer {
		return "", ErrInvalidToken
	}
	aud, err := t.Claims.GetAudience()
	if err != nil {
		return "", ErrInvalidToken
	}
	if len(aud) == 0 {
		return "", ErrInvalidToken
	}
	if aud[0] != audiences[0] {
		return "", ErrInvalidToken
	}

	userId, err := t.Claims.GetSubject()
	if err != nil {
		return "", ErrInvalidToken
	}

	return userId, nil
}

// AuthMiddleware returns JWT-based middleware to check authentication
// and set user id to context with contextKey
func AuthMiddleware(secret string, contextKey string) gin.HandlerFunc {
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
		userId, err := getUserIdFromToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
			return
		}
		c.Set(contextKey, userId)
		c.Next()
	}
}

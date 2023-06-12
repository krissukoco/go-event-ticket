package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-event-ticket/internal/auth"
	"github.com/krissukoco/go-event-ticket/internal/database"
	"github.com/krissukoco/go-event-ticket/internal/user"
)

const (
	defaultPort = 8080
)

type Server struct {
	Port int
	Pg   *sql.DB
}

func getPort() int {
	portStr, exists := os.LookupEnv("PORT")
	if !exists {
		return defaultPort
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return defaultPort
	}
	return port
}

func getJwtSecret() (string, error) {
	secret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		return "", fmt.Errorf("JWT_SECRET env not set")
	}
	return secret, nil
}

func migrateDb(db *sql.DB) error {
	return user.Migrate(db)
}

func main() {
	router := gin.Default()
	port := getPort()

	// Database
	db, err := database.NewPostgresDefault()
	if err != nil {
		panic(err)
	}
	if err = migrateDb(db); err != nil {
		panic(err)
	}

	v1 := router.Group("/api/v1")

	// Services
	userService := user.NewService(user.NewRepository(db))

	// Auth routes
	{
		secret, err := getJwtSecret()
		if err != nil {
			panic(err)
		}
		service := auth.NewService(userService, secret, 24*30)
		auth.RegisterHandlers(v1.Group("/auth"), service, auth.AuthMiddleware(service.UserIdFromToken))
	}
	// User routes
	{
		user.RegisterHandlers(v1.Group("/users"), userService)
	}

	router.Run(fmt.Sprintf(":%d", port))
}

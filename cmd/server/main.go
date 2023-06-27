package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-event-ticket/config"
	"github.com/krissukoco/go-event-ticket/internal/auth"
	"github.com/krissukoco/go-event-ticket/internal/database"
	"github.com/krissukoco/go-event-ticket/internal/user"
)

const (
	defaultPort   = 8080
	defaultJwtExp = 24 * 30
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

func main() {
	// Load config
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	cfg, err := config.Load(env)
	if err != nil {
		panic(err)
	}
	dbc := cfg.Database

	// Database
	db, err := database.NewPostgresGorm(dbc.Host, dbc.User, dbc.Password, dbc.DbName, dbc.Port, dbc.EnableSsl)
	if err != nil {
		panic(err)
	}

	// Server & Routers
	router := gin.Default()
	port := getPort()

	v1 := router.Group("/api/v1")

	// Services
	userService := user.NewService(user.NewRepository(db))

	// Auth routes
	authService := auth.NewService(userService, cfg.JwtSecret, defaultJwtExp)
	authMiddleware := auth.AuthMiddleware(authService.UserIdFromToken)
	auth.RegisterHandlers(v1.Group("/auth"), authService, authMiddleware)
	// User routes
	{
		user.RegisterHandlers(v1.Group("/users"), userService)
	}

	router.Run(fmt.Sprintf(":%d", port))
}

package main

import (
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
	authService := auth.NewService(userService, cfg.JwtSecret, defaultJwtExp)
	// Middlewares
	authMiddleware := auth.AuthMiddleware(cfg.JwtSecret, "userId")

	// Auth routes
	{
		authCtl := auth.NewController(authService, userService)
		authCtl.RegisterHandlers(v1.Group("/auth"))
	}
	// User routes
	{
		userCtl := user.NewController(userService, authMiddleware)
		userCtl.RegisterHandlers(v1.Group("/users"))
	}

	router.Run(fmt.Sprintf(":%d", port))
}

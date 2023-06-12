package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(group *gin.RouterGroup, service Service, authMiddleware gin.HandlerFunc) {
	group.POST("/login", login(service))
	group.POST("/register", register(service))

	// Private routes
	group.Use(authMiddleware)
	group.GET("/account", account(service))
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func login(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		r, err := service.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, r)
	}
}

func register(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		data := &registerData{
			Username: req.Username,
			Password: req.Password,
			Name:     req.Name,
			Location: req.Location,
		}
		r, err := service.Register(c.Request.Context(), data)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"id": r})
	}
}

func account(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userId")
		if userId == "" {
			// as the middleware auth is not used
			c.JSON(500, gin.H{"message": "internal server error"})
			return
		}
		r, err := service.Account(c.Request.Context(), userId)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, r)
	}
}

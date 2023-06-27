package user

import "github.com/gin-gonic/gin"

type controller struct {
	service        Service
	authMiddleware gin.HandlerFunc
}

func NewController(service Service, authMiddleware gin.HandlerFunc) *controller {
	return &controller{service, authMiddleware}
}

func (ctl *controller) RegisterHandlers(group *gin.RouterGroup) {
	group.GET("/users/me", ctl.authMiddleware, ctl.getAccount)
	group.GET("/users/:id", ctl.getUser)
}

func (ctl *controller) getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := ctl.service.GetById(id)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (ctl *controller) getAccount(c *gin.Context) {
	user, err := ctl.service.GetById(c.GetString("userId"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, user)
}

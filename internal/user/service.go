package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-event-ticket/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetById(id string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	ComparePassword(hashedPassword, password string) error
	HashPassword(password string) (string, error)
	Insert(username string, password string, name string, location string) (*models.User, error)
}

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrPasswordLength        = errors.New("password must be at least 8 characters")
)

type service struct {
	repo Repository
}

func RegisterHandlers(group *gin.RouterGroup, svc Service) {
	group.GET("/users/:id", getUser(svc))
}

func getUser(svc Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := svc.GetById(id)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, user)
	}
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetById(id string) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *service) GetByUsername(username string) (*models.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *service) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *service) validatePassword(pass string) error {
	if len(pass) < 8 {
		return ErrPasswordLength
	}
	return nil
}

func (s *service) Insert(username, password, name, location string) (*models.User, error) {
	// Check username already exists
	_, err := s.GetByUsername(username)
	if err == nil {
		return nil, ErrUsernameAlreadyExists
	}

	// Validate password
	if err = s.validatePassword(password); err != nil {
		return nil, err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Insert(username, hashedPassword, name, location)
	if err != nil {
		return nil, err
	}

	return user, nil
}

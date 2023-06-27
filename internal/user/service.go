package user

import (
	"context"
	"errors"

	"github.com/krissukoco/go-event-ticket/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetById(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	ComparePassword(hashedPassword, password string) error
	HashPassword(password string) (string, error)
	Insert(ctx context.Context, username string, password string, name string, location string) (*models.User, error)
}

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrPasswordLength        = errors.New("password must be at least 8 characters")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetById(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *service) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetByUsername(ctx, username)
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

func (s *service) Insert(ctx context.Context, username, password, name, location string) (*models.User, error) {
	// Check username already exists
	_, err := s.GetByUsername(ctx, username)
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

	user, err := s.repo.Insert(ctx, username, hashedPassword, name, location)
	if err != nil {
		return nil, err
	}

	return user, nil
}

package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/krissukoco/go-event-ticket/internal/models"
	"github.com/krissukoco/go-event-ticket/internal/user"
)

var (
	audiences                = []string{"github.com/krissukoco/go-event-ticket"}
	ErrInvalidToken          = errors.New("invalid token")
	ErrExpiredToken          = errors.New("expired token")
	ErrPasswordInvalid       = errors.New("password invalid")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

const (
	issuer = "github.com/krissukoco/go-event-ticket"
)

type loginData struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

type registerData struct {
	Username string
	Password string
	Name     string
	Location string
}

type Service interface {
	Login(ctx context.Context, username, password string) (*loginData, error)
	Register(ctx context.Context, in *registerData) (string, error)
	Account(ctx context.Context, userId string) (*models.User, error)
	UserIdFromToken(token string) (string, error)
}

type service struct {
	userService user.Service
	jwtSecret   string
	jwtExpHours int
}

func NewService(userService user.Service, jwtSecret string, expHours int) Service {
	return &service{userService, jwtSecret, expHours}
}

func (s *service) generateToken(userId string) (string, error) {
	now := time.Now()

	// Build claims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": userId,
		"aud": audiences,
		"exp": now.Add(time.Hour * time.Duration(s.jwtExpHours)).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"jti": "jwt_" + uuid.New().String(),
	})
	// Sign
	token, err := t.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) UserIdFromToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
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

func (s *service) Login(ctx context.Context, username, password string) (*loginData, error) {
	user, err := s.userService.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := s.userService.ComparePassword(user.Password, password); err != nil {
		return nil, ErrPasswordInvalid
	}

	token, err := s.generateToken(user.Id)
	if err != nil {
		return nil, err
	}

	return &loginData{User: user, Token: token}, nil
}

func (s *service) Register(ctx context.Context, in *registerData) (string, error) {
	// Check if username already exist
	_, err := s.userService.GetByUsername(in.Username)
	if err == nil {
		return "", ErrUsernameAlreadyExists
	}

	user, err := s.userService.Insert(in.Username, in.Password, in.Name, in.Location)
	if err != nil {
		return "", err
	}

	return user.Id, nil
}

func (s *service) Account(ctx context.Context, userId string) (*models.User, error) {
	user, err := s.userService.GetById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

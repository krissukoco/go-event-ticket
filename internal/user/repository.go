package user

import (
	"context"

	"github.com/krissukoco/go-event-ticket/internal/models"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Insert(ctx context.Context, username, password, name, location string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func newUserId() string {
	return "u_" + ulid.Make().String()
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&models.User{})

	return &repository{db}
}

func (r *repository) GetById(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Insert(ctx context.Context, username, password, name, location string) (*models.User, error) {
	id := newUserId()
	user := &models.User{
		Id:       id,
		Username: username,
		Password: password,
		Name:     name,
		Image:    "",
		Location: location,
	}

	tx := r.db.Create(user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return r.GetById(ctx, id)
}

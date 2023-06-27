package user

import (
	"github.com/krissukoco/go-event-ticket/internal/models"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Insert(username, password, name, location string) (*models.User, error)
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

func (r *repository) GetById(id string) (*models.User, error) {
	var u models.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Insert(username, password, name, location string) (*models.User, error) {
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

	return r.GetById(id)
}

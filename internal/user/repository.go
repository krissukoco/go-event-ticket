package user

import (
	"database/sql"
	"log"
	"time"
)

type Repository interface {
	GetById(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Insert(id, username, password, name, location string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func now() int64 {
	return time.Now().UnixMilli()
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) GetById(id string) (*User, error) {
	log.Println("GetById", id)
	var u User
	res := r.db.QueryRow("SELECT id, username, password, name, image, location, created_at, updated_at"+
		" FROM users WHERE id = $1", id)
	err := res.Scan(&u.Id, &u.Username, &u.Password, &u.Name, &u.Image, &u.Location, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetByUsername(username string) (*User, error) {
	var u User
	res := r.db.QueryRow("SELECT id, username, password, name, image, location, created_at, updated_at"+
		" FROM users WHERE username = $1", username)
	err := res.Scan(&u.Id, &u.Username, &u.Password, &u.Name, &u.Image, &u.Location, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

func (r *repository) Insert(id, username, password, name, location string) (*User, error) {
	ts := now()

	_, err := r.db.Exec("INSERT INTO users(id, username, password, name, image, location, created_at, updated_at)"+
		" VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING (id, username, password, name, image, location, created_at, updated_at)",
		id, username, password, name, "", location, ts, ts)

	if err != nil {
		return nil, err
	}

	return &User{
		Id:        id,
		Username:  username,
		Password:  password,
		Name:      name,
		Image:     "",
		Location:  location,
		CreatedAt: ts,
		UpdatedAt: ts,
	}, nil
}

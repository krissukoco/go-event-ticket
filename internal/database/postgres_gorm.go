package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func postgresDsnGorm(host, user, password, dbName string, port int, enableSsl bool) string {
	ssl := "disable"
	if enableSsl {
		ssl = "enable"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s Timezone=Asia/Jakarta", host, user, password, dbName, port, ssl)
}

func NewPostgresGorm(host, user, password, dbname string, port int, enableSsl bool) (*gorm.DB, error) {
	dsn := postgresDsnGorm(host, user, password, dbname, port, enableSsl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	defaultPostgresPort = 5432
)

var (
	ErrEnvHostNotSet  = errors.New("POSTGRES_HOST env not set")
	ErrEnvUserNotSet  = errors.New("POSTGRES_USER env not set")
	ErrEnvPassNotSet  = errors.New("POSTGRES_PASSWORD env not set")
	ErrEnvPortInvalid = errors.New("POSTGRES_PORT env invalid")
)

func postgresDsn(host, user, password, dbName string, disableSsl bool) string {
	sslMode := "enable"
	if disableSsl {
		sslMode = "disable"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", host, user, password, dbName, sslMode)
}

func NewPostgres(host, user, password, dbName string, port int, disableSsl bool) (*sql.DB, error) {
	db, err := sql.Open("postgres", postgresDsn(host, user, password, dbName, disableSsl))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// NewPostgresDefault returns a connection to a local postgres database
// with default credentials on environment variables.
func NewPostgresDefault() (*sql.DB, error) {
	host, exists := os.LookupEnv("POSTGRES_HOST")
	if !exists {
		return nil, ErrEnvHostNotSet
	}
	user, exists := os.LookupEnv("POSTGRES_USER")
	if !exists {
		return nil, ErrEnvUserNotSet
	}
	password, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if !exists {
		return nil, ErrEnvPassNotSet
	}
	dbName, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		dbName = "go_event_ticket"
	}
	port := defaultPostgresPort
	portStr, exists := os.LookupEnv("POSTGRES_PORT")
	if exists {
		convPort, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, ErrEnvPortInvalid
		}
		port = convPort
	}
	disableSsl := true
	sslEnv, exists := os.LookupEnv("POSTGRES_SSL")
	if exists {
		if sslEnv == "enable" {
			disableSsl = false
		}
	}

	return NewPostgres(host, user, password, dbName, port, disableSsl)
}

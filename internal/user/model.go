package user

import "database/sql"

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Location  string `json:"location"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			image TEXT,
			location TEXT,
			created_at BIGINT NOT NULL,
			updated_at BIGINT NOT NULL
		);
	`)
	return err
}

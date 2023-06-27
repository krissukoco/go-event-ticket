package models

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

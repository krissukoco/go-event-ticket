package event

import "database/sql"

// Event represents an event entity in the database
type Event struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	// StartTime and EndTime are in milliseconds
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	// Start is a string in the format "YYYY-MM-DD HH:MM"
	Start string `json:"start"`
	// End is a string in the format "YYYY-MM-DD HH:MM"
	End string `json:"end"`
	// Timezone is a string in the format "Asia/Jakarta"
	Timezone   string `json:"timezone"`
	LocationId string `json:"location_id"`
	CreatorId  string `json:"creator_id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// EventTicket represents ticket that can be bought for an event
type EventTicket struct {
	Id string `json:"id"`
	// EventId is the id of the event that this ticket is for
	EventId string `json:"event_id"`
	// Name is the name of the ticket (e.g. "Early Bird", "VIP", etc.)
	Name string `json:"name"`
	// Description is the description of the ticket (e.g. "20% Cheaper!", "Free T-Shirt", etc.)
	Description string `json:"description"`
	// Price is the price of the ticket in IDR
	Price int64 `json:"price"`
	// Quantity is the number of tickets available for this type
	Quantity int `json:"quantity"`
	// StartAt and EndAt are in milliseconds
	StartAt int64 `json:"start_at"`
	EndAt   int64 `json:"end_at"`
	// CreatedAt and UpdatedAt are in milliseconds
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func MigrateEvents(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		start_time BIGINT NOT NULL,
		end_time BIGINT NOT NULL,
		start TEXT NOT NULL,
		end TEXT NOT NULL,
		timezone TEXT NOT NULL,
		location_id TEXT NOT NULL,
		creator_id TEXT NOT NULL,
		created_at BIGINT NOT NULL,
		updated_at BIGINT NOT NULL
	);`
	_, err := db.Exec(q)
	return err
}

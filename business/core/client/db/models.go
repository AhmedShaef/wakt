package db

import (
	"time"
)

// Client represent the structure we need for moving data
// between the app and the database.
type Client struct {
	ID          string    `db:"client_id"`
	Name        string    `db:"name"`
	Wid         string    `db:"wid"`
	Notes       string    `db:"notes"`
	DateUpdated time.Time `db:"date_updated"`
}

package db

import (
	"time"
)

// Group represent the structure we need for moving data
// between the app and the database.
type Group struct {
	ID          string    `db:"group_id"`
	Name        string    `db:"name"`
	Wid         string    `db:"wid"`
	DateUpdated time.Time `db:"date_updated"`
}

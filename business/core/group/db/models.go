package db

import (
	"time"
)

// Group represent the structure we need for moving data
// between the app and the database.
type Group struct {
	ID          string    `db:"group_id"`
	Name        string    `db:"name"`
	WID         string    `db:"wid"`
	UID         string    `db:"uid"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

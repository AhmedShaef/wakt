package db

import "time"

// Tag represent the structure we need for moving data
// between the app and the database.
type Tag struct {
	ID          string    `db:"tag_id"`
	Name        string    `db:"name"`
	WID         string    `db:"wid"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

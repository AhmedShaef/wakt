package db

import (
	"time"
)

// Team represent the structure we need for moving data
// between the app and the database.
type Team struct {
	ID          string    `db:"team_id"`
	PID         string    `db:"pid"`
	UID         string    `db:"uid"`
	WID         string    `db:"wid"`
	Manager     bool      `db:"manager"`
	Rate        float64   `db:"rate"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

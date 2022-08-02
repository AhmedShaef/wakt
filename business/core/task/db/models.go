package db

import (
	"time"
)

// Task represent the structure we need for moving data
// between the app and the database.
type Task struct {
	ID             string        `db:"task_id"`
	Name           string        `db:"name"`
	Pid            string        `db:"pid"`
	Wid            string        `db:"wid"`
	UID            string        `db:"uid"`
	Estimated      time.Duration `db:"estimated_seconds"`
	Active         bool          `db:"active"`
	DateCreated    time.Time     `db:"date_created"`
	DateUpdated    time.Time     `db:"date_updated"`
	TrackedSeconds time.Duration `db:"tracked_seconds"`
}

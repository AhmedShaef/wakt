package db

import (
	"time"
)

// Task represent the structure we need for moving data
// between the app and the database.
type Task struct {
	ID               string        `db:"task_id"`
	Name             string        `db:"name"`
	Pid              string        `db:"pid"`
	Wid              string        `db:"wid"`
	Uid              string        `db:"uid"`
	EstimatedSeconds time.Duration `db:"estimated_seconds"`
	Active           bool          `db:"active"`
	DateUpdated      time.Time     `db:"date_updated"`
	TrackedSeconds   time.Duration `db:"tracked_seconds"`
}

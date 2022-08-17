package db

import (
	"time"

	"github.com/lib/pq"
)

// TimeEntry represent the structure we need for moving data
// between the app and the database.
type TimeEntry struct {
	ID          string         `db:"time_entry_id"`
	Description string         `db:"description"`
	UID         string         `db:"uid"`
	WID         string         `db:"wid"`
	PID         string         `db:"pid"`
	TID         string         `db:"tid"`
	Billable    bool           `db:"billable"`
	Start       time.Time      `db:"start"`
	Stop        time.Time      `db:"stop"`
	Duration    time.Duration  `db:"duration"`
	CreatedWith string         `db:"created_with"`
	Tags        pq.StringArray `db:"tags"`
	DurOnly     bool           `db:"dur_only"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
}

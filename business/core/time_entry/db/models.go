package db

import (
	"time"
)

// TimeEntrie represent the structure we need for moving data
// between the app and the database.
type TimeEntrie struct {
	ID          string        `db:"time_entrie_id"`
	Description string        `db:"description"`
	Uid         string        `db:"uid"`
	Wid         string        `db:"wid"`
	Pid         string        `db:"pid"`
	Tid         string        `db:"tid"`
	Billable    bool          `db:"billable"`
	Start       time.Time     `db:"start"`
	Stop        time.Time     `db:"stop"`
	Duration    time.Duration `db:"duration"`
	CreatedWith string        `db:"created_with"`
	Tags        []string      `db:"tags"`
	DurOnly     bool          `db:"dur_only"`
	DateCreated time.Time     `db:"date_created"`
	DateUpdated time.Time     `db:"date_updated"`
}

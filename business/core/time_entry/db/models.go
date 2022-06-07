package db

import (
	"github.com/lib/pq"
	"time"
)

// TimeEntry represent the structure we need for moving data
// between the app and the database.
type TimeEntry struct {
	ID          string         `db:"time_entry_id"`
	Description string         `db:"description" default:"no description"`
	Uid         string         `db:"uid"`
	Wid         string         `db:"wid"`
	Pid         string         `db:"pid" default:""`
	Tid         string         `db:"tid" default:""`
	Billable    bool           `db:"billable" default:"false"`
	Start       time.Time      `db:"start"`
	Stop        time.Time      `db:"stop" default:""`
	Duration    time.Duration  `db:"duration" default:"-1"`
	CreatedWith string         `db:"created_with"`
	Tags        pq.StringArray `db:"tags" default:"[]"`
	DurOnly     bool           `db:"dur_only" default:"false"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
}

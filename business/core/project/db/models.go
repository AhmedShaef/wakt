package db

import "time"

// Project represent the structure we need for moving data
// between the app and the database.
type Project struct {
	ID             string        `db:"project_id"`
	Name           string        `db:"name"`
	Wid            string        `db:"wid"`
	Cid            string        `db:"cid" default:""`
	Active         bool          `db:"active" default:"true"`
	IsPrivate      bool          `db:"is_private" default:"true"`
	Billable       bool          `db:"billable" default:"true"`
	AutoEstimates  bool          `db:"auto_estimates" default:"false"`
	EstimatedHours time.Duration `db:"estimated_hours" default:""`
	DateCreated    time.Time     `db:"date_created"`
	DateUpdated    time.Time     `db:"date_updated"`
	Rate           float32       `db:"rate" default:""`
	HexColor       string        `db:"hex_color" default:""`
}

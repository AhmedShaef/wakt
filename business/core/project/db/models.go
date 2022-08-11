package db

import "time"

// Project represent the structure we need for moving data
// between the app and the database.
type Project struct {
	ID             string        `db:"project_id"`
	Name           string        `db:"name"`
	WID            string        `db:"wid"`
	Cid            string        `db:"cid"`
	UID            string        `db:"uid"`
	Active         bool          `db:"active"`
	IsPrivate      bool          `db:"is_private"`
	Billable       bool          `db:"billable"`
	AutoEstimates  bool          `db:"auto_estimates"`
	EstimatedHours time.Duration `db:"estimated_hours"`
	DateCreated    time.Time     `db:"date_created"`
	DateUpdated    time.Time     `db:"date_updated"`
	Rate           float32       `db:"rate"`
	HexColor       string        `db:"hex_color"`
}

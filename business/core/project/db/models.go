package db

import "time"

// Project represent the structure we need for moving data
// between the app and the database.
type Project struct {
	ID             string        `db:"project_id"`
	Name           string        `db:"name"`
	Wid            string        `db:"wid"`
	Cid            string        `db:"cid"`
	Active         bool          `db:"active"`
	IsPrivate      bool          `db:"is_private"`
	Template       bool          `db:"template"`
	TemplateID     string        `db:"template_id"`
	Billable       bool          `db:"billable"`
	AutoEstimates  bool          `db:"auto_estimates"`
	EstimatedHours time.Duration `db:"estimated_hours"`
	DateUpdated    time.Time     `db:"date_updated"`
	Color          int           `db:"color"`
	Rate           float32       `db:"rate"`
	DateCreated    time.Time     `db:"date_created"`
}

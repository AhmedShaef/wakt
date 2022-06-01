package db

import (
	"time"
)

// ProjectUser represent the structure we need for moving data
// between the app and the database.
type ProjectUser struct {
	ID          string    `db:"project_user_id"`
	Pid         string    `db:"pid"`
	Uid         string    `db:"uid"`
	Wid         string    `db:"wid"`
	Manager     bool      `db:"manager"`
	Rate        float64   `db:"rate"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

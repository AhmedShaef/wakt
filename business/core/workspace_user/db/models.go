package db

import "time"

// WorkspaceUser represent the structure we need for moving data
// between the app and the database.
type WorkspaceUser struct {
	ID          string    `db:"workspace_user_id"`
	Uid         string    `db:"uid"`
	Wid         string    `db:"wid"`
	Admin       bool      `db:"admin"`
	Active      bool      `db:"active"`
	InviteKey   string    `db:"invite_key"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

package db

// WorkspaceUser represent the structure we need for moving data
// between the app and the database.
type WorkspaceUser struct {
	ID        string `db:"workspace_user_id"`
	Uid       string `db:"uid"`
	Wid       string `db:"wid"`
	Active    bool   `db:"active"`
	InviteKey string `db:"invite_key"`
}

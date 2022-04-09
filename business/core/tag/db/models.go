package db

// Tag represent the structure we need for moving data
// between the app and the database.
type Tag struct {
	ID   string `db:"tag_id"`
	Name string `db:"name"`
	Wid  string `db:"wid"`
}

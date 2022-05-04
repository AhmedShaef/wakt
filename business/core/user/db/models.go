package db

import (
	"github.com/lib/pq"
	"time"
)

// User represent the structure we need for moving data
// between the app and the database.
type User struct {
	ID              string         `db:"user_id"`
	DefaultWid      string         `db:"default_wid"`
	Email           string         `db:"email"`
	PasswordHash    []byte         `db:"password_hash"`
	Roles           pq.StringArray `db:"roles"`
	FullName        string         `db:"full_name"`
	TimeOfDayFormat string         `db:"time_of_day_format"`
	DateFormat      string         `db:"date_format"`
	BeginningOfWeek int            `db:"beginning_of_week"`
	Language        string         `db:"language"`
	ImageURL        string         `db:"image_url"`
	DateCreated     time.Time      `db:"date_created"`
	DateUpdated     time.Time      `db:"date_updated"`
	TimeZone        string         `db:"timezone"`
	Invitation      pq.StringArray `db:"invitation"`
	DurationFormat  string         `db:"duration_format"`
}

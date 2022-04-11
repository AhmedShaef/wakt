package db

import (
	"time"
)

// User represent the structure we need for moving data
// between the app and the database.
type User struct {
	ID                     string    `db:"user_id"`
	APIToken               string    `db:"api_token"` //TODO generate tokens
	DefaultWid             string    `db:"default_wid"`
	Email                  string    `db:"email"`
	PasswordHash           []byte    `db:"password_hash"`
	Roles                  []string  `db:"roles"`
	FullName               string    `db:"full_name"`
	JqueryTimeOfDayFormat  string    `db:"jquery_time_of_day_format"`
	JqueryDateFormat       string    `db:"jquery_date_format"`
	TimeOfDayFormat        string    `db:"time_of_day_format"`
	DateFormat             string    `db:"date_format"`
	StoreStartAndStopTime  bool      `db:"store_start_and_stop_time"`
	BeginningOfWeek        int       `db:"beginning_of_week"`
	Language               string    `db:"language"`
	ImageURL               string    `db:"image_url"`
	SidebarPiechart        bool      `db:"sidebar_piechart"`
	DateCreated            time.Time `db:"date_created"`
	DateUpdated            time.Time `db:"date_updated"`
	RecordTimeline         bool      `db:"record_timeline"`
	ShouldUpgrade          bool      `db:"should_upgrade"`
	SendProductEmails      bool      `db:"send_product_emails"`
	SendWeeklyReport       bool      `db:"send_weekly_report"`
	SendTimerNotifications bool      `db:"send_timer_notifications"`
	OpenidEnabled          bool      `db:"openid_enabled"`
	TimeZone               string    `db:"timezone"`
	Invitation             []string  `db:"invitation"`
	DurationFormat         string    `db:"duration_format"`
}

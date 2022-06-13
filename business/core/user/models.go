package user

import (
	"github.com/AhmedShaef/wakt/business/core/user/db"
	"time"
	"unsafe"
)

// User represents an individual user.
type User struct {
	ID              string    `json:"id"`
	DefaultWid      string    `json:"default_wid"`
	Email           string    `json:"email"`
	PasswordHash    []byte    `json:"password_hash"`
	FullName        string    `json:"full_name"`
	TimeOfDayFormat string    `json:"time_of_day_format"`
	DateFormat      string    `json:"date_format"`
	BeginningOfWeek int       `json:"beginning_of_week"`
	Language        string    `json:"language"`
	ImageURL        string    `json:"image_url"`
	DateCreated     time.Time `json:"date_created"`
	DateUpdated     time.Time `json:"date_updated"`
	TimeZone        string    `json:"timezone"`
	Invitation      []string  `json:"invitation"`
	DurationFormat  string    `json:"duration_format"`
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

// UpdateUser defines what information may be provided to modify an existing
// user. All fields are optional so user can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateUser struct {
	DefaultWid      *string    `json:"default_wid"`
	Email           *string    `json:"email" validate:"omitempty,email"`
	FullName        *string    `json:"full_name"`
	TimeOfDayFormat *string    `json:"time_of_day_format" validate:"omitempty,eq=H:mm|eq=h:mm"`
	DateFormat      *string    `json:"date_format" validate:"omitempty,eq=YYYY-MM-DD|eq=DD.MM.YYYY|eq=DD-MM-YYYY|eq=MM/DD/YYYY|eq=DD/MM/YYYY|eq=MM-DD-YYYY"`
	BeginningOfWeek *int       `json:"beginning_of_week" valiate:"omitempty,min=0|max=6"`
	Language        *string    `json:"language"`
	DateCreated     *time.Time `json:"date_created"`
	DateUpdated     *time.Time `json:"date_updated"`
	TimeZone        *string    `json:"timezone"`
	Invitation      []string   `json:"invitation"`
	DurationFormat  *string    `json:"duration_format"`
}

//ChangePassword defines what information may be provided to change an existing
// user's password.
type ChangePassword struct {
	OldPassword string  `json:"old_password" validate:"required,min=6,max=64"`
	Password    *string `json:"Password" validate:"required,min=6,max=64"`
}

// UpdateImage defines what information may be provided to update an existing
// user's image.
type UpdateImage struct {
	ImageName string `json:"image_name" validate:"required,omitempty"`
}

// =============================================================================

func toUser(dbUser db.User) User {
	pu := (*User)(unsafe.Pointer(&dbUser))
	return *pu
}

func toUsersSlice(dbUser []db.User) []User {
	users := make([]User, len(dbUser))
	for i, dbusr := range dbUser {
		users[i] = toUser(dbusr)
	}
	return users
}

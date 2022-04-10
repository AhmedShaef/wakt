package time_entry

import (
	"github.com/AhmedShaef/wakt/business/core/time_entry/db"
	"time"
	"unsafe"
)

// TimeEntrie represents an individual time_entry.
type TimeEntrie struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Uid         string        `json:"uid"`
	Wid         string        `json:"wid"`
	Pid         string        `json:"pid"`
	Tid         string        `json:"tid"`
	Billable    bool          `json:"billable"`
	Start       time.Time     `json:"start"`
	Stop        time.Time     `json:"stop"`
	Duration    time.Duration `json:"duration"`
	CreatedWith string        `json:"created_with"`
	Tags        []string      `json:"tags"`
	DurOnly     bool          `json:"dur_only"`
	DateCreated time.Time     `json:"date_created"`
	DateUpdated time.Time     `json:"date_updated"`
}

// NewTimeEntrie contains information needed to create a new time_entry.
type NewTimeEntrie struct {
	Description string        `json:"description"`
	Wid         string        `json:"wid"`
	Pid         string        `json:"pid"`
	Tid         string        `json:"tid"`
	Billable    bool          `json:"billable"`
	Start       time.Time     `json:"start" validate:"required"`
	Stop        time.Time     `json:"stop"`
	Duration    time.Duration `json:"duration" validate:"required"`
	CreatedWith string        `json:"created_with" validate:"required"`
	Tags        []string      `json:"tags"`
	DurOnly     bool          `json:"dur_only"`
}

//StartTimeEntrie contains information needed to start a new time_entry.
type StartTimeEntrie struct {
	Description string   `json:"description"`
	Wid         string   `json:"wid"`
	Pid         string   `json:"pid"`
	Tid         string   `json:"tid"`
	Billable    bool     `json:"billable"`
	CreatedWith string   `json:"created_with" validate:"required"`
	Tags        []string `json:"tags"`
	DurOnly     bool     `json:"dur_only"`
}

// UpdateTimeEntrie defines what information may be provided to modify an existing
// time_entry. All fields are optional so time_entry can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateTimeEntrie struct {
	Description *string        `json:"description"`
	Pid         *string        `json:"pid"`
	Tid         *string        `json:"tid"`
	Billable    *bool          `json:"billable"`
	Start       *time.Time     `json:"start"`
	Stop        *time.Time     `json:"stop"`
	Duration    *time.Duration `json:"duration"`
	Tags        *[]string      `json:"tags"`
	DurOnly     *bool          `json:"dur_only"`
}

//UpdateTimeEntrieTags contains information needed to update bulk of time_entry tags.
type UpdateTimeEntrieTags struct {
	Tags    *[]string `json:"tags" validate:"required"`
	TagMode string    `json:"tag_mode" validate:"required eq=add|eq=remove"`
}

// =============================================================================

func toTimeEntrie(dbTimeEntrie db.TimeEntrie) TimeEntrie {
	pu := (*TimeEntrie)(unsafe.Pointer(&dbTimeEntrie))
	return *pu
}

func toTimeEntrieSlice(dbTimeEntrie []db.TimeEntrie) []TimeEntrie {
	TimeEntrie := make([]TimeEntrie, len(dbTimeEntrie))
	for i, dbtimEntri := range dbTimeEntrie {
		TimeEntrie[i] = toTimeEntrie(dbtimEntri)
	}
	return TimeEntrie
}

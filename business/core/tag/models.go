package tag

import (
	"github.com/AhmedShaef/wakt/business/core/tag/db"
	"time"
	"unsafe"
)

// Tag represents an individual tag.
type Tag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Wid         string    `json:"wid"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewTag contains information needed to create a new tag.
type NewTag struct {
	Name string `json:"name" validate:"required"`
	Wid  string `json:"wid" validate:"required"`
}

// UpdateTag defines what information may be provided to modify an existing
// tag. All fields are optional so tag can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateTag struct {
	Name *string `json:"name"`
}

// =============================================================================

func toTag(dbTag db.Tag) Tag {
	pu := (*Tag)(unsafe.Pointer(&dbTag))
	return *pu
}

func toTagsSlice(dbTag []db.Tag) []Tag {
	tags := make([]Tag, len(dbTag))
	for i, dbtag := range dbTag {
		tags[i] = toTag(dbtag)
	}
	return tags
}

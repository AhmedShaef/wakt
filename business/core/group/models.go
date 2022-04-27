package group

import (
	"github.com/AhmedShaef/wakt/business/core/group/db"
	"time"
	"unsafe"
)

// Group represents an individual Group.
type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Wid         string    `json:"wid"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewGroup contains information needed to create a new Group.
type NewGroup struct {
	Name string `json:"name" validate:"required"`
	Wid  string `json:"wid" validate:"required"`
}

// UpdateGroup defines what information may be provided to modify an existing
// group. All fields are optional so group can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateGroup struct {
	Name *string `json:"name"`
}

// =============================================================================

func toGroup(dbGroup db.Group) Group {
	pu := (*Group)(unsafe.Pointer(&dbGroup))
	return *pu
}

func toGroupSlice(dbGroup []db.Group) []Group {
	groups := make([]Group, len(dbGroup))
	for i, dbgrop := range dbGroup {
		groups[i] = toGroup(dbgrop)
	}
	return groups
}

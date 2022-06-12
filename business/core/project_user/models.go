package project_user

import (
	"github.com/AhmedShaef/wakt/business/core/project_user/db"
	"time"
	"unsafe"
)

// ProjectUser represents an individual ProjectUser.
type ProjectUser struct {
	ID          string    `json:"project_user_id"`
	Pid         string    `json:"pid"`
	Uid         string    `json:"uid"`
	Wid         string    `json:"wid"`
	Manager     bool      `json:"manager"`
	Rate        float64   `json:"rate"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewProjectUser contains information needed to create a new ProjectUser.
type NewProjectUser struct {
	Pid     string  `json:"pid" validate:"required"`
	Uid     string  `json:"uid" validate:"required"`
	Wid     string  `json:"wid"`
	Manager bool    `json:"manager"`
	Rate    float64 `json:"rate"`
	Puis    string  `json:"puis"`
}

// UpdateProjectUser defines what information may be provided to modify an existing
// project user. All fields are optional so project users can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateProjectUser struct {
	Rate    *float64 `json:"rate"`
	Manager *bool    `json:"manager"`
}

// =============================================================================

func toProjectUser(dbProjectUser db.ProjectUser) ProjectUser {
	pu := (*ProjectUser)(unsafe.Pointer(&dbProjectUser))
	return *pu
}

func toProjectUserSlice(dbProjectUser []db.ProjectUser) []ProjectUser {
	projectUsers := make([]ProjectUser, len(dbProjectUser))
	for i, dbProjectUser := range dbProjectUser {
		projectUsers[i] = toProjectUser(dbProjectUser)
	}
	return projectUsers
}

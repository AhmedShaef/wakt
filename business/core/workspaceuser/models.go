package workspaceuser

import (
	"time"
	"unsafe"

	"github.com/AhmedShaef/wakt/business/core/workspaceuser/db"
)

// WorkspaceUser represents a workspace user
type WorkspaceUser struct {
	ID          string    `json:"id"`
	UID         string    `json:"uid"`
	Wid         string    `json:"wid"`
	Admin       bool      `json:"admin"`
	Active      bool      `json:"active"`
	InviteKey   string    `json:"invite_key"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// InviteUsers contains information to invite users to a workspace.
type InviteUsers struct {
	Emails    []string `json:"emails" validate:"required"`
	InviterID string   `json:"inviter_id" validate:"required"`
}

// UpdateWorkspaceUser defines what information may be provided to modify an existing
// client. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateWorkspaceUser struct {
	Active *bool `json:"active"`
	Admin  *bool `json:"admin"`
}

//==============================================================================

func toWorkspaceUser(dbWorkspaceUser db.WorkspaceUser) WorkspaceUser {
	wu := (*WorkspaceUser)(unsafe.Pointer(&dbWorkspaceUser))
	return *wu
}

func toWorkspaceUserSlice(dbWorkspaceUsers []db.WorkspaceUser) []WorkspaceUser {
	workspaceUsers := make([]WorkspaceUser, len(dbWorkspaceUsers))
	for i, dbWorkspaceUser := range dbWorkspaceUsers {
		workspaceUsers[i] = toWorkspaceUser(dbWorkspaceUser)
	}
	return workspaceUsers
}

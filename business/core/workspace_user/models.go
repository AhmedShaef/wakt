package workspace_user

import (
	"github.com/AhmedShaef/wakt/business/core/workspace_user/db"
	"unsafe"
)

// WorkspaceUser represents a workspace user
type WorkspaceUser struct {
	Id        string `json:"id"`
	Uid       string `json:"uid"`
	Wid       string `json:"wid"`
	Active    bool   `json:"active"`
	InviteKey string `json:"invite_key"`
}

// InviteUsers contains information to invite users to a workspace.
type InviteUsers struct {
	Emails []string `json:"emails" validate:"required"`
}

// UpdateWorkspaceUser defines what information may be provided to modify an existing
// client. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields ,so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types ,but we make exceptions around
// marshalling/unmarshalling.
type UpdateWorkspaceUser struct {
	Active *bool `json:"active"`
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

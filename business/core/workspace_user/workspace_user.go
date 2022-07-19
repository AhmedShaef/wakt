// Package workspace_user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package workspace_user

import (
	"context"
	"errors"
	"fmt"
	"time"

	users "github.com/AhmedShaef/wakt/business/core/user/db"
	"github.com/AhmedShaef/wakt/business/core/workspace_user/db"
	send "github.com/AhmedShaef/wakt/business/send/smtp"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("user not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for user access.
type Core struct {
	store     db.Store
	userStore users.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store:     db.NewStore(log, sqlxDB),
		userStore: users.NewStore(log, sqlxDB),
	}
}

// Create add a user to a workspace.
func (c Core) Create(ctx context.Context, workspaceID string, userID string, now time.Time) (WorkspaceUser, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return WorkspaceUser{}, ErrInvalidID
	}

	if err := validate.CheckID(userID); err != nil {
		return WorkspaceUser{}, ErrInvalidID
	}

	dbWorkspaceUser := db.WorkspaceUser{
		ID:          validate.GenerateID(),
		UID:         userID,
		Wid:         workspaceID,
		Admin:       true,
		Active:      true,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.store.Invite(ctx, dbWorkspaceUser); err != nil {
		return WorkspaceUser{}, fmt.Errorf("invite: %w", err)
	}

	return toWorkspaceUser(dbWorkspaceUser), nil
}

// InviteUser invites a user to a workspace.
func (c Core) InviteUser(ctx context.Context, workspaceID string, ni InviteUsers, now time.Time) ([]WorkspaceUser, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return []WorkspaceUser{}, ErrInvalidID
	}

	if err := validate.Check(ni); err != nil {
		return []WorkspaceUser{}, fmt.Errorf("validating data: %w", err)
	}

	var workspaceUsers []db.WorkspaceUser

	// limit the invitation by 10 new user and check
	//if the input emails is in valid email syntax.
	var emails []string
	for i, v := range ni.Emails {
		if validate.CheckEmail(v) && i < 10 {
			emails = append(emails, v)
		}
	}

	for _, v := range emails {
		var userID string

		usr, err := c.userStore.QueryByEmail(ctx, v)
		if err != nil {
			if errors.Is(err, database.ErrDBNotFound) {
				userID = validate.GenerateID()
			} else {
				return []WorkspaceUser{}, fmt.Errorf("query: %w", err)
			}

		} else {
			userID = usr.ID
			dbUser := users.User{
				Invitation: append(usr.Invitation, workspaceID),
			}
			dbUser.DateUpdated = now
			if err := c.userStore.Update(ctx, dbUser); err != nil {
				return []WorkspaceUser{}, fmt.Errorf("udpate: %w", err)
			}
		}

		dbWorkspaceUser := db.WorkspaceUser{
			ID:          validate.GenerateID(),
			UID:         userID,
			Wid:         workspaceID,
			Admin:       false,
			Active:      false,
			InviteKey:   "?invite_key=" + "68hh6542-" + userID + "-9fo5d7l5-" + workspaceID + "-9y6d4y2b5r6o",
			DateCreated: now,
			DateUpdated: now,
		}
		workspaceUsers = append(workspaceUsers, dbWorkspaceUser)

		if err := send.Email("example@example.com", v, "You have been invited to join WAKT!", "www.example.com/signup/"+dbWorkspaceUser.InviteKey); err != nil {
			return []WorkspaceUser{}, fmt.Errorf("send email: %w", err)
		}
		if err := c.store.Invite(ctx, dbWorkspaceUser); err != nil {
			return []WorkspaceUser{}, fmt.Errorf("invite: %w", err)
		}

	}

	return toWorkspaceUserSlice(workspaceUsers), nil
}

// Update replaces a workspace user document in the database.
func (c Core) Update(ctx context.Context, workspaceUserID string, uwu UpdateWorkspaceUser, now time.Time) error {
	if err := validate.CheckID(workspaceUserID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uwu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbWorkspaceUser, err := c.store.QueryByID(ctx, workspaceUserID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating workspaceUser workspaceUserID[%s]: %w", workspaceUserID, err)
	}

	if uwu.Active != nil {
		dbWorkspaceUser.Active = *uwu.Active
	}
	if uwu.Admin != nil {
		dbWorkspaceUser.Admin = *uwu.Admin
	}
	dbWorkspaceUser.DateUpdated = now

	if err := c.store.Update(ctx, dbWorkspaceUser); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a workspace user from the database.
func (c Core) Delete(ctx context.Context, workspaceUserID string) error {
	if err := validate.CheckID(workspaceUserID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, workspaceUserID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

//Query retrieves a list of existing workspace users from the database.
func (c Core) Query(ctx context.Context, workspaceID string, pageNumber int, rowsPerPage int) ([]WorkspaceUser, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return []WorkspaceUser{}, ErrInvalidID
	}
	dbWorkspaceUsers, err := c.store.Query(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toWorkspaceUserSlice(dbWorkspaceUsers), nil
}

// QueryByID gets the specified workspaceUser from the database.
func (c Core) QueryByID(ctx context.Context, workspaceUserID string) (WorkspaceUser, error) {
	if err := validate.CheckID(workspaceUserID); err != nil {
		return WorkspaceUser{}, ErrInvalidID
	}

	dbWorkspaceUser, err := c.store.QueryByID(ctx, workspaceUserID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return WorkspaceUser{}, ErrNotFound
		}
		return WorkspaceUser{}, fmt.Errorf("query: %w", err)
	}

	return toWorkspaceUser(dbWorkspaceUser), nil
}

// QueryByuIDwID gets the specified workspaceUser from the database.
func (c Core) QueryByuIDwID(ctx context.Context, workspaceID, userID string) (WorkspaceUser, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return WorkspaceUser{}, ErrInvalidID
	}
	if err := validate.CheckID(userID); err != nil {
		return WorkspaceUser{}, ErrInvalidID
	}

	dbWorkspaceUser, err := c.store.QueryByuIDwID(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return WorkspaceUser{}, ErrNotFound
		}
		return WorkspaceUser{}, fmt.Errorf("query: %w", err)
	}

	return toWorkspaceUser(dbWorkspaceUser), nil
}

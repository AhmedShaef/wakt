// Package workspace_user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package workspace_user

import (
	"context"
	"errors"
	"fmt"
	users "github.com/AhmedShaef/wakt/business/core/user/db"
	"github.com/AhmedShaef/wakt/business/core/workspace_user/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/smtp"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("user not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for user access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// InviteUser invites a user to a workspace.
func (c Core) InviteUser(ctx context.Context, workspaceID string, ni InviteUsers, cfg smtp.Config) ([]WorkspaceUser, error) {
	if err := validate.Check(ni); err != nil {
		return []WorkspaceUser{}, fmt.Errorf("validating data: %w", err)
	}

	var workspaceUsers []db.WorkspaceUser

	for _, i := range ni.Emails {
		var userID string
		usr, err := users.Store{}.QueryByEmail(ctx, i)
		if err != nil {
			println(err)
		}

		if usr.ID != "" {
			userID = usr.ID
		} else {
			userID = validate.GenerateID()
		}

		dbWorkspaceUser := db.WorkspaceUser{
			ID:        validate.GenerateID(),
			Uid:       userID,
			Wid:       workspaceID,
			InviteKey: validate.GenerateID(),
		}
		workspaceUsers = append(workspaceUsers, dbWorkspaceUser)

		if err := smtp.SendEmail(cfg, i, "You have been invited to join WAKT!", ""); err != nil {
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
func (c Core) Query(ctx context.Context, userID string, pageNumber int, rowsPerPage int) ([]WorkspaceUser, error) {
	dbWorkspaceUsers, err := c.store.Query(ctx, userID, pageNumber, rowsPerPage)
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

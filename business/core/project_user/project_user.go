// Package project_user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package project_user

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/project_user/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
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

// Create inserts a new project user into the database.
func (c Core) Create(ctx context.Context, npu NewProjectUser, now time.Time) ([]ProjectUser, error) {
	if err := validate.Check(npu); err != nil {
		return []ProjectUser{}, fmt.Errorf("validating data: %w", err)
	}

	uids := strings.Split(npu.Uid, ",")

	var ProjectUsers []db.ProjectUser

	for _, uid := range uids {

		dbprojectuser := db.ProjectUser{
			Rate: 0,
		}
		dbprojectuser = db.ProjectUser{
			ID:          validate.GenerateID(),
			Pid:         npu.Pid,
			Uid:         uid,
			Wid:         npu.Wid,
			Manager:     npu.Manager,
			Rate:        npu.Rate,
			DateCreated: now,
			DateUpdated: now,
		}

		if err := c.store.Create(ctx, dbprojectuser); err != nil {
			return []ProjectUser{}, fmt.Errorf("create: %w", err)
		}
		ProjectUsers = append(ProjectUsers, dbprojectuser)
	}

	return toProjectUserSlice(ProjectUsers), nil
}

// Update replaces a project user document in the database.
func (c Core) Update(ctx context.Context, projectUserID string, upu UpdateProjectUser, now time.Time) error {
	if err := validate.CheckID(projectUserID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(upu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbProjectUser, err := c.store.QueryByID(ctx, projectUserID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating projectUser projectUserID[%s]: %w", projectUserID, err)
	}

	if upu.Rate != nil {
		dbProjectUser.Rate = *upu.Rate
	}
	if upu.Manager != nil {
		dbProjectUser.Manager = *upu.Manager
	}
	dbProjectUser.DateUpdated = now

	if err := c.store.Update(ctx, dbProjectUser); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a project user from the database.
func (c Core) Delete(ctx context.Context, projectUserID string) error {
	if err := validate.CheckID(projectUserID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, projectUserID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryWorkspaceProjectUsers retrieves a list of existing project user from the database.
func (c Core) QueryWorkspaceProjectUsers(ctx context.Context, WorkspaceID string, pageNumber, rowsPerPage int) ([]ProjectUser, error) {
	if err := validate.CheckID(WorkspaceID); err != nil {
		return []ProjectUser{}, ErrInvalidID
	}
	dbProjectUser, err := c.store.QueryWorkspaceProjectUsers(ctx, WorkspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toProjectUserSlice(dbProjectUser), nil
}

// QueryByID gets the specified project user from the database.
func (c Core) QueryByID(ctx context.Context, projectUserID string) (ProjectUser, error) {
	if err := validate.CheckID(projectUserID); err != nil {
		return ProjectUser{}, ErrInvalidID
	}

	dbProjectUser, err := c.store.QueryByID(ctx, projectUserID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ProjectUser{}, ErrNotFound
		}
		return ProjectUser{}, fmt.Errorf("query: %w", err)
	}

	return toProjectUser(dbProjectUser), nil
}

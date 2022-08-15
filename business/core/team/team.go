// Package team provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package team

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AhmedShaef/wakt/business/core/team/db"
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
	store db.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create inserts a new project user into the database.
func (c Core) Create(ctx context.Context, npu NewTeam, now time.Time) ([]Team, error) {
	if err := validate.Check(npu); err != nil {
		return []Team{}, fmt.Errorf("validating data: %w", err)
	}

	uids := strings.Split(npu.UID, ",")

	var Teams []db.Team

	for _, uid := range uids {

		dbprojectuser := db.Team{
			Rate: 0,
		}
		dbprojectuser = db.Team{
			ID:          validate.GenerateID(),
			PID:         npu.Pid,
			UID:         uid,
			Wid:         npu.Wid,
			Manager:     npu.Manager,
			Rate:        npu.Rate,
			DateCreated: now,
			DateUpdated: now,
		}

		if err := c.store.Create(ctx, dbprojectuser); err != nil {
			return []Team{}, fmt.Errorf("create: %w", err)
		}
		Teams = append(Teams, dbprojectuser)
	}

	return toTeamSlice(Teams), nil
}

// Update replaces a project user document in the database.
func (c Core) Update(ctx context.Context, projectUserID string, upu UpdateTeam, now time.Time) error {
	if err := validate.CheckID(projectUserID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(upu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbTeam, err := c.store.QueryByID(ctx, projectUserID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating projectUser projectUserID[%s]: %w", projectUserID, err)
	}

	if upu.Rate != nil {
		dbTeam.Rate = *upu.Rate
	}
	if upu.Manager != nil {
		dbTeam.Manager = *upu.Manager
	}
	dbTeam.DateUpdated = now

	if err := c.store.Update(ctx, dbTeam); err != nil {
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

// QueryWorkspaceTeams retrieves a list of existing project user from the database.
func (c Core) QueryWorkspaceTeams(ctx context.Context, WorkspaceID string, pageNumber, rowsPerPage int) ([]Team, error) {
	if err := validate.CheckID(WorkspaceID); err != nil {
		return []Team{}, ErrInvalidID
	}
	dbTeam, err := c.store.QueryWorkspaceTeams(ctx, WorkspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toTeamSlice(dbTeam), nil
}

// QueryByID gets the specified project user from the database.
func (c Core) QueryByID(ctx context.Context, projectUserID string) (Team, error) {
	if err := validate.CheckID(projectUserID); err != nil {
		return Team{}, ErrInvalidID
	}

	dbTeam, err := c.store.QueryByID(ctx, projectUserID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Team{}, ErrNotFound
		}
		return Team{}, fmt.Errorf("query: %w", err)
	}

	return toTeam(dbTeam), nil
}

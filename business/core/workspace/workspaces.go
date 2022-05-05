// Package workspace provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package workspace

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/workspace/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
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

// Create inserts a new Workspace into the database.
func (c Core) Create(ctx context.Context, nw NewWorkspace, now time.Time) (Workspace, error) {
	if err := validate.Check(nw); err != nil {
		return Workspace{}, fmt.Errorf("validating data: %w", err)
	}

	dbworkspace := db.Workspace{
		ID:          validate.GenerateID(),
		Name:        nw.Name,
		Uid:         nw.Uid,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.store.Create(ctx, dbworkspace); err != nil {
		return Workspace{}, fmt.Errorf("create: %w", err)
	}
	return toWorkspace(dbworkspace), nil
}

// Update replaces a workspace document in the database.
func (c Core) Update(ctx context.Context, workspaceID string, uw UpdateWorkspace, now time.Time) error {
	if err := validate.CheckID(workspaceID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uw); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbWorkspace, err := c.store.QueryByID(ctx, workspaceID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating workspace workspaceID[%s]: %w", workspaceID, err)
	}

	if uw.Name != nil {
		dbWorkspace.Name = *uw.Name
	}
	if uw.DefaultHourlyRate != nil {
		dbWorkspace.DefaultHourlyRate = *uw.DefaultHourlyRate
	}
	if uw.DefaultCurrency != nil {
		dbWorkspace.DefaultCurrency = *uw.DefaultCurrency
	}
	if uw.OnlyAdminMayCreateProjects != nil {
		dbWorkspace.OnlyAdminMayCreateProjects = *uw.OnlyAdminMayCreateProjects
	}
	if uw.OnlyAdminSeeBillableRates != nil {
		dbWorkspace.OnlyAdminSeeBillableRates = *uw.OnlyAdminSeeBillableRates
	}
	if uw.OnlyAdminSeeTeamDashboard != nil {
		dbWorkspace.OnlyAdminSeeTeamDashboard = *uw.OnlyAdminSeeTeamDashboard
	}
	if uw.ProjectBillableByDefault != nil {
		dbWorkspace.ProjectBillableByDefault = *uw.ProjectBillableByDefault
	}
	if uw.Rounding != nil {
		dbWorkspace.Rounding = *uw.Rounding
	}
	if uw.RoundingMinutes != nil {
		dbWorkspace.RoundingMinutes = *uw.RoundingMinutes
	}
	dbWorkspace.DateUpdated = now

	if err := c.store.Update(ctx, dbWorkspace); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// UpdateLogo replaces a user document in the database.
func (c Core) UpdateLogo(ctx context.Context, workspaceID string, uw UpdateWorkspace, now time.Time) error {
	if err := validate.CheckID(workspaceID); err != nil {
		return ErrInvalidID
	}

	dbuser, err := c.store.QueryByID(ctx, workspaceID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating user workspaceID[%s]: %w", workspaceID, err)
	}
	if uw.LogoURL != "" {
		dbuser.LogoURL = "app/tooling/uploader/assets" + uw.LogoURL
	}
	dbuser.DateUpdated = now

	if err := c.store.Update(ctx, dbuser); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

//Query retrieves a list of existing workspace from the database.
func (c Core) Query(ctx context.Context, userID string, pageNumber int, rowsPerPage int) ([]Workspace, error) {
	dbWorkspace, err := c.store.Query(ctx, userID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toWorkspaceSlice(dbWorkspace), nil
}

// QueryByID gets the specified workspace from the database.
func (c Core) QueryByID(ctx context.Context, workspaceID string) (Workspace, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return Workspace{}, ErrInvalidID
	}

	dbWorkspace, err := c.store.QueryByID(ctx, workspaceID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Workspace{}, ErrNotFound
		}
		return Workspace{}, fmt.Errorf("query: %w", err)
	}

	return toWorkspace(dbWorkspace), nil
}

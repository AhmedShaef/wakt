// Package group provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package group

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AhmedShaef/wakt/business/core/group/db"
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

// Create inserts a new group into the database.
func (c Core) Create(ctx context.Context, userID string, ng NewGroup, now time.Time) (Group, error) {
	if err := validate.CheckID(userID); err != nil {
		return Group{}, ErrInvalidID
	}
	if err := validate.Check(ng); err != nil {
		return Group{}, fmt.Errorf("validating data: %w", err)
	}

	dbgrop := db.Group{
		ID:          validate.GenerateID(),
		Name:        ng.Name,
		Wid:         ng.WID,
		UID:         userID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.store.Create(ctx, dbgrop); err != nil {
		return Group{}, fmt.Errorf("create: %w", err)
	}

	return toGroup(dbgrop), nil
}

// Update replaces a group document in the database.
func (c Core) Update(ctx context.Context, groupID string, ug UpdateGroup, now time.Time) error {
	if err := validate.CheckID(groupID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(ug); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbgrop, err := c.store.QueryByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating group groupID[%s]: %w", groupID, err)
	}

	if ug.Name != nil {
		dbgrop.Name = *ug.Name
	}
	dbgrop.DateUpdated = now

	if err := c.store.Update(ctx, dbgrop); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a group from the database.
func (c Core) Delete(ctx context.Context, groupID string) error {
	if err := validate.CheckID(groupID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, groupID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByID gets the specified group from the database.
func (c Core) QueryByID(ctx context.Context, groupID string) (Group, error) {
	if err := validate.CheckID(groupID); err != nil {
		return Group{}, ErrInvalidID
	}

	dbgrop, err := c.store.QueryByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Group{}, ErrNotFound
		}
		return Group{}, fmt.Errorf("query: %w", err)
	}

	return toGroup(dbgrop), nil
}

// QueryWorkspaceGroups retrieves a list of existing workspace from the database.
func (c Core) QueryWorkspaceGroups(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Group, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return []Group{}, ErrInvalidID
	}

	dbGroups, err := c.store.QueryWorkspaceGroups(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return toGroupSlice(dbGroups), nil
}

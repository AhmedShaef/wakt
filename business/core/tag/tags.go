// Package tag provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package tag

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/tag/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("tag not found")
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

// Create inserts a new tag into the database.
func (c Core) Create(ctx context.Context, nt NewTag, now time.Time) (Tag, error) {
	if err := validate.Check(nt); err != nil {
		return Tag{}, fmt.Errorf("validating data: %w", err)
	}

	dbtg := db.Tag{
		ID:          validate.GenerateID(),
		Name:        nt.Name,
		WID:         nt.Wid,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.store.Create(ctx, dbtg); err != nil {
		return Tag{}, fmt.Errorf("create: %w", err)
	}

	return toTag(dbtg), nil
}

// Update replaces a tag document in the database.
func (c Core) Update(ctx context.Context, tagID string, ut UpdateTag, now time.Time) error {
	if err := validate.CheckID(tagID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(ut); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbtg, err := c.store.QueryByID(ctx, tagID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating tag tagID[%s]: %w", tagID, err)
	}

	if ut.Name != nil {
		dbtg.Name = *ut.Name
	}
	dbtg.DateUpdated = now
	if err := c.store.Update(ctx, dbtg); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a tag from the database.
func (c Core) Delete(ctx context.Context, tagID string) error {
	if err := validate.CheckID(tagID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, tagID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByID gets the specified tag from the database.
func (c Core) QueryByID(ctx context.Context, tagID string) (Tag, error) {
	if err := validate.CheckID(tagID); err != nil {
		return Tag{}, ErrInvalidID
	}

	dbtg, err := c.store.QueryByID(ctx, tagID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Tag{}, ErrNotFound
		}
		return Tag{}, fmt.Errorf("query: %w", err)
	}

	return toTag(dbtg), nil
}

// QueryWorkspaceTags retrieves a list of existing workspace from the database.
func (c Core) QueryWorkspaceTags(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Tag, error) {
	if err := validate.CheckID(workspaceID); err != nil {
		return []Tag{}, ErrInvalidID
	}
	dbTags, err := c.store.QueryWorkspaceTags(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return toTagsSlice(dbTags), nil
}

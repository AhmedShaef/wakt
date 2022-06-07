// Package tasks provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.

package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/task/db"
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

// Create inserts a new task into the database.
func (c Core) Create(ctx context.Context, userID string, nt NewTask, now time.Time) (Task, error) {
	if err := validate.Check(nt); err != nil {
		return Task{}, fmt.Errorf("validating data: %w", err)
	}

	nameInProject := c.store.QueryUnique(ctx, nt.Name, "pid", nt.Pid)
	if nameInProject != "" {
		return Task{}, fmt.Errorf("project name is not unique for workspace")
	}

	dbtask := db.Task{
		Uid:              "00000000-0000-0000-0000-000000000000",
		EstimatedSeconds: 0,
		Active:           true,
		TrackedSeconds:   0,
	}
	dbtask = db.Task{
		ID:               validate.GenerateID(),
		Name:             nt.Name,
		Pid:              nt.Pid,
		Wid:              nt.Wid,
		Uid:              userID,
		EstimatedSeconds: nt.EstimatedSeconds,
		Active:           nt.Active,
		DateCreated:      now,
		DateUpdated:      now,
		TrackedSeconds:   nt.TrackedSeconds,
	}

	if err := c.store.Create(ctx, dbtask); err != nil {
		return Task{}, fmt.Errorf("create: %w", err)
	}

	return toTask(dbtask), nil
}

// Update replaces a task document in the database.
func (c Core) Update(ctx context.Context, taskID string, uc UpdateTask, now time.Time) error {
	if err := validate.CheckID(taskID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uc); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbtask, err := c.store.QueryByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating task taskID[%s]: %w", taskID, err)
	}

	if uc.Name != nil {
		dbtask.Name = *uc.Name
	}
	if uc.EstimatedSeconds != nil {
		dbtask.EstimatedSeconds = *uc.EstimatedSeconds
	}
	if uc.Active != nil {
		dbtask.Active = *uc.Active
	}
	if uc.TrackedSeconds != nil {
		dbtask.TrackedSeconds = *uc.TrackedSeconds
	}
	dbtask.DateUpdated = now

	if err := c.store.Update(ctx, dbtask); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a task from the database.
func (c Core) Delete(ctx context.Context, taskID string) error {
	if err := validate.CheckID(taskID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, taskID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByID gets the specified task from the database.
func (c Core) QueryByID(ctx context.Context, taskID string) (Task, error) {
	if err := validate.CheckID(taskID); err != nil {
		return Task{}, ErrInvalidID
	}

	dbtsk, err := c.store.QueryByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Task{}, ErrNotFound
		}
		return Task{}, fmt.Errorf("query: %w", err)
	}

	return toTask(dbtsk), nil
}

//QueryProjectTasks retrieves a list of existing projects from the database.
func (c Core) QueryProjectTasks(ctx context.Context, projectID string, pageNumber, rowsPerPage int) ([]Task, error) {
	dbprojects, err := c.store.QueryProjectTasks(ctx, projectID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toTasksSlice(dbprojects), nil
}

// QueryWorkspaceTasks retrieves a list of existing workspace from the database.
func (c Core) QueryWorkspaceTasks(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Task, error) {
	dbTasks, err := c.store.QueryWorkspaceTasks(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return toTasksSlice(dbTasks), nil
}

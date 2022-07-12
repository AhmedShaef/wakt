package time_entry

import (
	"context"
	"errors"
	"fmt"

	dbp "github.com/AhmedShaef/wakt/business/core/project/db"
	dbt "github.com/AhmedShaef/wakt/business/core/task/db"
	"github.com/AhmedShaef/wakt/business/core/time_entry/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/util"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"

	"time"

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

// Create inserts a new time entry into the database.
func (c Core) Create(ctx context.Context, nt NewTimeEntry, userID string, now time.Time) (TimeEntry, error) {
	if err := validate.CheckID(userID); err != nil {
		return TimeEntry{}, ErrInvalidID
	}

	if err := validate.Check(nt); err != nil {
		return TimeEntry{}, fmt.Errorf("validating data: %w", err)
	}

	stop := nt.Start.Add(nt.Duration)

	// set values from request
	dbTimeEntry := db.TimeEntry{
		ID:          validate.GenerateID(),
		Description: nt.Description,
		UID:         userID,
		Wid:         nt.Wid,
		Pid:         nt.Pid,
		Tid:         nt.Tid,
		Billable:    nt.Billable,
		Start:       nt.Start,
		Stop:        stop,
		Duration:    nt.Duration,
		CreatedWith: nt.CreatedWith,
		Tags:        nt.Tags,
		DurOnly:     nt.DurOnly,
		DateCreated: now,
		DateUpdated: now,
	}
	if dbTimeEntry.Pid == "" {
		dbTimeEntry.Pid = "00000000-0000-0000-0000-000000000000"
	}
	if dbTimeEntry.Tid == "" {
		dbTimeEntry.Tid = "00000000-0000-0000-0000-000000000000"
	}
	if dbTimeEntry.Tags == nil {
		dbTimeEntry.Tags = []string{}
	}

	if err := c.store.Create(ctx, dbTimeEntry); err != nil {
		return TimeEntry{}, fmt.Errorf("create: %w", err)
	}

	if err := c.SyncTaskTime(ctx, dbTimeEntry.Tid, now); err != nil {
		return TimeEntry{}, fmt.Errorf("sync task time: %w", err)
	}

	if err := c.SyncProjectTime(ctx, dbTimeEntry.Pid, now); err != nil {
		return TimeEntry{}, fmt.Errorf("sync project time: %w", err)
	}

	return toTimeEntry(dbTimeEntry), nil
}

// Start inserts a new time entry into the database.
func (c Core) Start(ctx context.Context, st StartTimeEntry, userID string, now time.Time) (TimeEntry, error) {
	if err := validate.CheckID(userID); err != nil {
		return TimeEntry{}, ErrInvalidID
	}

	if err := validate.Check(st); err != nil {
		return TimeEntry{}, fmt.Errorf("validating data: %w", err)
	}

	// set values from request
	dbTimeEntry := db.TimeEntry{
		ID:          validate.GenerateID(),
		Description: st.Description,
		UID:         userID,
		Wid:         st.Wid,
		Pid:         st.Pid,
		Tid:         st.Tid,
		Billable:    st.Billable,
		Start:       now,
		Stop:        time.Time{},
		Duration:    -1,
		CreatedWith: st.CreatedWith,
		Tags:        st.Tags,
		DurOnly:     st.DurOnly,
		DateCreated: now,
		DateUpdated: now,
	}

	if dbTimeEntry.Pid == "" {
		dbTimeEntry.Pid = "00000000-0000-0000-0000-000000000000"
	}
	if dbTimeEntry.Tid == "" {
		dbTimeEntry.Tid = "00000000-0000-0000-0000-000000000000"
	}
	if dbTimeEntry.Tags == nil {
		dbTimeEntry.Tags = []string{}
	}

	if err := c.store.Create(ctx, dbTimeEntry); err != nil {
		return TimeEntry{}, fmt.Errorf("create: %w", err)
	}

	return toTimeEntry(dbTimeEntry), nil
}

// Stop replaces a time_entry document in the database.
func (c Core) Stop(ctx context.Context, TimeEntryID string, now time.Time) (TimeEntry, error) {
	if err := validate.CheckID(TimeEntryID); err != nil {
		return TimeEntry{}, ErrInvalidID
	}

	dbTimeEntry, err := c.store.QueryByID(ctx, TimeEntryID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return TimeEntry{}, ErrNotFound
		}
		return TimeEntry{}, fmt.Errorf("stopping time_entry time_entryID[%s]: %w", TimeEntryID, err)
	}

	dbTimeEntry.Stop = now
	dbTimeEntry.Duration = dbTimeEntry.Stop.Sub(dbTimeEntry.Start)
	dbTimeEntry.DateUpdated = now

	if err := c.store.Update(ctx, dbTimeEntry); err != nil {
		return TimeEntry{}, fmt.Errorf("stop: %w", err)
	}

	if err := c.SyncTaskTime(ctx, dbTimeEntry.Tid, now); err != nil {
		return TimeEntry{}, fmt.Errorf("sync task time: %w", err)
	}

	if err := c.SyncProjectTime(ctx, dbTimeEntry.Pid, now); err != nil {
		return TimeEntry{}, fmt.Errorf("sync project time: %w", err)
	}

	return toTimeEntry(dbTimeEntry), nil
}

// Update replaces a time_entry document in the database.
func (c Core) Update(ctx context.Context, TimeEntryID string, ut UpdateTimeEntry, now time.Time) error {
	if err := validate.CheckID(TimeEntryID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(ut); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbTimEntry, err := c.store.QueryByID(ctx, TimeEntryID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating time_entry time_entryID[%s]: %w", TimeEntryID, err)
	}

	if ut.Description != nil {
		dbTimEntry.Description = *ut.Description
	}
	if ut.Billable != nil {
		dbTimEntry.Billable = *ut.Billable
	}
	if ut.Start != nil {
		dbTimEntry.Start = *ut.Start
		dbTimEntry.Duration = dbTimEntry.Stop.Sub(dbTimEntry.Start)
	}
	if ut.Stop != nil {
		dbTimEntry.Stop = *ut.Stop
		dbTimEntry.Duration = dbTimEntry.Stop.Sub(dbTimEntry.Start)
	}
	if ut.CreatedWith != nil {
		dbTimEntry.CreatedWith = *ut.CreatedWith
	}
	if ut.Tags != nil {
		dbTimEntry.Tags = ut.Tags
	}
	if ut.DurOnly != nil {
		dbTimEntry.DurOnly = *ut.DurOnly
	}
	dbTimEntry.DateUpdated = now

	if err := c.store.Update(ctx, dbTimEntry); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a time_entry from the database.
func (c Core) Delete(ctx context.Context, timeEntryID string) error {
	if err := validate.CheckID(timeEntryID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, timeEntryID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByID gets the specified time_entry from the database.
func (c Core) QueryByID(ctx context.Context, timeEntryID string) (TimeEntry, error) {
	if err := validate.CheckID(timeEntryID); err != nil {
		return TimeEntry{}, ErrInvalidID
	}

	dbTimEntry, err := c.store.QueryByID(ctx, timeEntryID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return TimeEntry{}, ErrNotFound
		}
		return TimeEntry{}, fmt.Errorf("query: %w", err)
	}
	return toTimeEntry(dbTimEntry), nil
}

// SyncProjectTime sync the specified project time from the database.
func (c Core) SyncProjectTime(ctx context.Context, projectID string, now time.Time) error {
	if err := validate.CheckID(projectID); err != nil {
		return ErrInvalidID
	}

	ProjctTime, err := c.store.QueryProjectTime(ctx, projectID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("sync: %w", err)
	}

	dbproject := dbp.Project{
		ID:             projectID,
		EstimatedHours: ProjctTime.Duration / 60 / 60,
		DateUpdated:    now,
	}

	if err = c.store.UpdateProjectTime(ctx, dbproject); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

// SyncTaskTime sync the specified task time from the database.
func (c Core) SyncTaskTime(ctx context.Context, taskID string, now time.Time) error {
	if err := validate.CheckID(taskID); err != nil {
		return ErrInvalidID
	}

	taskTime, err := c.store.QueryTaskTime(ctx, taskID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("sync: %w", err)
	}
	dbTask := dbt.Task{
		ID:             taskID,
		TrackedSeconds: taskTime.Duration,
		DateUpdated:    now,
	}
	if err = c.store.UpdateTaskTime(ctx, dbTask); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}
	return nil
}

//QueryRunning retrieves a list of existing time entry from the database.
func (c Core) QueryRunning(ctx context.Context, userID string, pageNumber int, rowsPerPage int) ([]TimeEntry, error) {
	if err := validate.CheckID(userID); err != nil {
		return []TimeEntry{}, ErrInvalidID
	}

	dbTimeEntry, err := c.store.QueryRunning(ctx, userID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toTimeEntrySlice(dbTimeEntry), nil
}

//QueryRange retrieves a list of existing time entry from the database.
func (c Core) QueryRange(ctx context.Context, userID string, pageNumber, rowsPerPage int, start, end time.Time) ([]TimeEntry, error) {
	if err := validate.CheckID(userID); err != nil {
		return []TimeEntry{}, ErrInvalidID
	}

	dbTimeEntry, err := c.store.QueryRange(ctx, userID, pageNumber, rowsPerPage, start, end)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toTimeEntrySlice(dbTimeEntry), nil
}

// UpdateTags replaces a time_entry document in the database.
func (c Core) UpdateTags(ctx context.Context, TimeEntryID string, ut UpdateTimeEntryTags, now time.Time) error {
	if err := validate.CheckID(TimeEntryID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(ut); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbTimEntry, err := c.store.QueryByID(ctx, TimeEntryID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating time_entry time_entryID[%s]: %w", TimeEntryID, err)
	}

	if ut.TagMode == "add" {
		dbTimEntry.Tags = util.Add(dbTimEntry.Tags, ut.Tags)
	} else if ut.TagMode == "remove" {
		dbTimEntry.Tags = util.Remove(dbTimEntry.Tags, ut.Tags)
	}
	dbTimEntry.DateUpdated = now
	if err := c.store.Update(ctx, dbTimEntry); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

//QueryDash retrieves a list of existing time entry from the database.
func (c Core) QueryDash(ctx context.Context, UserID string) ([]TimeEntry, error) {
	if err := validate.CheckID(UserID); err != nil {
		return []TimeEntry{}, ErrInvalidID
	}

	var dbTimeEntrys []db.TimeEntry

	dbMostActive, err := c.store.QueryMostActive(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	dbActiveity, err := c.store.QueryActivity(ctx, UserID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	dbTimeEntrys = append(dbTimeEntrys, dbMostActive...)

	dbTimeEntrys = append(dbTimeEntrys, dbActiveity...)

	return toTimeEntrySlice(dbTimeEntrys), nil
}

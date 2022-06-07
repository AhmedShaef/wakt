// Package db contains TimeEntry related CRUD functionality.
package db

import (
	"context"
	"fmt"
	dbp "github.com/AhmedShaef/wakt/business/core/project/db"
	"github.com/AhmedShaef/wakt/business/core/task/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// Store manages the set of APIs for user access.
type Store struct {
	log          *zap.SugaredLogger
	tr           database.Transactor
	db           sqlx.ExtContext
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		tr:  db,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		return fn(s.db)
	}
	return database.WithinTran(ctx, s.log, s.tr, fn)
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// Create adds a TimeEntry to the database. It returns the created TimeEntry with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, te TimeEntry) error {
	const q = `
	INSERT INTO time_entries
		(time_entry_id, description, uid, wid, pid, tid, billable, start, stop, duration, tags, created_with, dur_only, date_created, date_updated)
	VALUES
		(:time_entry_id, :description, :uid, :wid, :pid, :tid, :billable, :start, :stop, :duration, :tags, :created_with, :dur_only, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, te); err != nil {
		return fmt.Errorf("inserting time_entry: %w", err)
	}

	return nil
}

// Update modifies data about a TimeEntry. It will error if the specified ID is
// invalid or does not reference an existing TimeEntry.
func (s Store) Update(ctx context.Context, te TimeEntry) error {
	const q = `
	UPDATE
		time_entries
	SET
		"description" = :description,
		"uid" = :uid,
		"pid" = :pid,
		"tid" = :tid,
        "billable" = :billable,
        "start" = :start,
        "stop" = :stop,
        "duration" = :duration,
        "dur_only" = :dur_only,
		"tags" = :tags,
		"date_updated" = :date_updated
	WHERE
		time_entry_id = :time_entry_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, te); err != nil {
		return fmt.Errorf("updating time_entry time_entryID[%s]: %w", te.ID, err)
	}

	return nil
}

// Delete removes the TimeEntry identified by a given ID.
func (s Store) Delete(ctx context.Context, timeEntryID string) error {
	data := struct {
		TimeEntryID string `db:"time_entry_id"`
	}{
		TimeEntryID: timeEntryID,
	}

	const q = `
	DELETE FROM
		time_entries
	WHERE
		time_entry_id = :time_entry_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting time_entry time_entryID[%s]: %w", timeEntryID, err)
	}

	return nil
}

// QueryByID finds the time_entry identified by a given ID.
func (s Store) QueryByID(ctx context.Context, timeEntryID string) (TimeEntry, error) {
	data := struct {
		TimeEntryID string `db:"time_entry_id"`
	}{
		TimeEntryID: timeEntryID,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
		time_entry_id = :time_entry_id`

	var tim TimeEntry
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &tim); err != nil {
		return TimeEntry{}, fmt.Errorf("selecting time_entry time_entryID[%q]: %w", timeEntryID, err)
	}

	return tim, nil
}

// QueryRunning gets all TimeEntry from the database.
func (s Store) QueryRunning(ctx context.Context, userID string, pageNumber int, rowsPerPage int) ([]TimeEntry, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		UserID      string `db:"user_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		UserID:      userID,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
		duration < 0
		AND uid = :user_id
	ORDER BY
		time_entry_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tims []TimeEntry
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryRange gets all TimeEntry from the database.
func (s Store) QueryRange(ctx context.Context, userID string, pageNumber, rowsPerPage int, start, end time.Time) ([]TimeEntry, error) {
	data := struct {
		Offset      int       `db:"offset"`
		RowsPerPage int       `db:"rows_per_page"`
		Start       time.Time `db:"start"`
		End         time.Time `db:"end"`
		UserID      string    `db:"user_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		Start:       start,
		End:         end,
		UserID:      userID,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
		date_created >= :start AND date_created <= :end
		AND uid = :user_id
	ORDER BY
		time_entry_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tims []TimeEntry
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryMostActive user in all TimeEntry from the database.
func (s Store) QueryMostActive(ctx context.Context, userID string) ([]TimeEntry, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}
	const q = `
	SELECT
		uid, duration
	FROM
		time_entries
	WHERE
		stop >= now() - INTERVAL '1 week'
		AND uid = :user_id
	ORDER BY
		duration
	LIMIT 5`

	var tims []TimeEntry
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryActivity gets last TimeEntry from the database.
func (s Store) QueryActivity(ctx context.Context, userID string) ([]TimeEntry, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}
	const q = `
	SELECT
		uid, pid, duration, description, stop, tid
	FROM
		time_entries
	WHERE
		uid = :user_id
	ORDER BY
		start DESC
	LIMIT 20`

	var tims []TimeEntry
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryProjectTime sync specified project time from the database.
func (s Store) QueryProjectTime(ctx context.Context, projectID string) (TimeEntry, error) {
	data := struct {
		ProjectID string `db:"project_id"`
	}{
		ProjectID: projectID,
	}

	const q = `
	SELECT
		SUM(duration) AS duration
	FROM
		time_entries
	WHERE 
		pid = :project_id`

	var ps TimeEntry
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &ps); err != nil {
		return TimeEntry{}, fmt.Errorf("selecting tracked time projectID[%q]: %w", projectID, err)
	}

	return ps, nil
}

// UpdateProjectTime modifies data about a TimeEntry. It will error if the specified ID is
// invalid or does not reference an existing TimeEntry.
func (s Store) UpdateProjectTime(ctx context.Context, data dbp.Project) error {

	const q2 = `
	UPDATE
		projects
	SET
		"estimated_hours"= :estimated_hours,
		"date_updated" = :date_updated
	WHERE
		"project_id" = :project_id AND "auto_estimates" = true`

	if err := database.NamedExecContext(ctx, s.log, s.db, q2, data); err != nil {
		return fmt.Errorf("updating tracked_second projectID[%s]: %w", data.ID, err)
	}

	return nil
}

// QueryTaskTime sync the specified task time from the database.
func (s Store) QueryTaskTime(ctx context.Context, taskID string) (TimeEntry, error) {
	data := struct {
		TaskID string `db:"task_id"`
	}{
		TaskID: taskID,
	}

	const q = `
	SELECT
		SUM(duration) AS duration
	FROM
		time_entries
	WHERE 
		tid = :task_id`

	var ts TimeEntry
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &ts); err != nil {
		return TimeEntry{}, fmt.Errorf("selecting tracked time taskID[%q]: %w", taskID, err)
	}

	return ts, nil
}

// UpdateTaskTime modifies data about a TimeEntry. It will error if the specified ID is
// invalid or does not reference an existing TimeEntry.
func (s Store) UpdateTaskTime(ctx context.Context, data db.Task) error {

	const q2 = `
	UPDATE
		tasks
	SET
		"tracked_seconds"= :tracked_seconds,
		"date_updated"=:date_updated
	WHERE
		"task_id"= :task_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q2, data); err != nil {
		return fmt.Errorf("updating tracked_second taskID[%s]: %w", data.ID, err)
	}

	return nil
}

// Package db contains TimeEntrie related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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

// Create adds a TimeEntrie to the database. It returns the created TimeEntrie with
// fields like ID and DateCreated populated.
func (s Store) Create(ctx context.Context, te TimeEntrie) error {
	const q = `
	INSERT INTO time_entries
		(time_entrie_id, description, wid, pid, tid, billable, start, stop, duration, created_with, tags, dur_only, date_created, date_updated)
	VALUES
		(:time_entrie_id, :description, :wid, :pid, :tid, :billable, :start, :stop, :duration, :created_with, :tags, :dur_only, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, te); err != nil {
		return fmt.Errorf("inserting time_entry: %w", err)
	}

	return nil
}

// Update modifies data about a TimeEntrie. It will error if the specified ID is
// invalid or does not reference an existing TimeEntrie.
func (s Store) Update(ctx context.Context, te TimeEntrie) error {
	const q = `
	UPDATE
		time_entries
	SET
		"description" = :description,
		"pid" = :pid,
		"tid" = :tid,
        "billable" = :billable,
        "start" = :start,
        "stop" = :stop,
        "duration" = :duration,
        "tags" = :tags,
        "dur_only" = :dur_only,
		"date_updated" = :date_updated
	WHERE
		time_entrie_id = :time_entrie_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, te); err != nil {
		return fmt.Errorf("updating time_entry time_entrieID[%s]: %w", te.ID, err)
	}

	return nil
}

// Delete removes the TimeEntrie identified by a given ID.
func (s Store) Delete(ctx context.Context, timeEntrieID string) error {
	data := struct {
		TimeEntrieID string `db:"time_entrie_id"`
	}{
		TimeEntrieID: timeEntrieID,
	}

	const q = `
	DELETE FROM
		time_entries
	WHERE
		time_entrie_id = :time_entrie_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting time_entry time_entrieID[%s]: %w", timeEntrieID, err)
	}

	return nil
}

// Query gets all TimeEntrie from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]TimeEntrie, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	ORDER BY
		time_entrie_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tims []TimeEntrie
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryByID finds the time_entry identified by a given ID.
func (s Store) QueryByID(ctx context.Context, timeEntrieID string) (TimeEntrie, error) {
	data := struct {
		TimeEntrieID string `db:"time_entrie_id"`
	}{
		TimeEntrieID: timeEntrieID,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
		time_entrie_id = :time_entrie_id`

	var tim TimeEntrie
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &tim); err != nil {
		return TimeEntrie{}, fmt.Errorf("selecting time_entry time_entrieID[%q]: %w", timeEntrieID, err)
	}

	return tim, nil
}

// QueryRunning gets all TimeEntrie from the database.
func (s Store) QueryRunning(ctx context.Context, pageNumber int, rowsPerPage int) ([]TimeEntrie, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
		duration < 0
	ORDER BY
		time_entrie_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tims []TimeEntrie
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryRange gets all TimeEntrie from the database.
func (s Store) QueryRange(ctx context.Context, pageNumber, rowsPerPage int, start, end string) ([]TimeEntrie, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		Start       string `db:"start"`
		End         string `db:"end"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		Start:       start,
		End:         end,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
			date_created >= :start AND date_created <= :end
	ORDER BY
		time_entrie_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tims []TimeEntrie
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryBulkIDs gets all TimeEntrie from the database.
func (s Store) QueryBulkIDs(ctx context.Context, timeEntrieID []string) ([]TimeEntrie, error) {
	data := struct {
		TimeEntrieID []string `db:"time_entrie_id"`
	}{
		TimeEntrieID: timeEntrieID,
	}

	const q = `
	SELECT
		*
	FROM
		time_entries
	WHERE 
		time_entrie_id in :time_entrie_id`

	var tims []TimeEntrie
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryMostActive user in all TimeEntrie from the database.
func (s Store) QueryMostActive(ctx context.Context) ([]TimeEntrie, error) {
	data := struct{}{}
	const q = `
	SELECT
		uid, duration
	FROM
		time_entries
	WHERE
		stop >= now() - INTERVAL '1 week'
	ORDER BY
		duration
	LIMIT 5`

	var tims []TimeEntrie
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

// QueryActivity gets last TimeEntrie from the database.
func (s Store) QueryActivity(ctx context.Context) ([]TimeEntrie, error) {
	data := struct{}{}
	const q = `
	SELECT
		uid, pid, duration, description, stop, tid
	FROM
		time_entries
	ORDER BY
		start DESC
	LIMIT 20`

	var tims []TimeEntrie
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tims); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return tims, nil
}

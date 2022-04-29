// Package db contains task related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/jmoiron/sqlx"
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

// Create inserts a new task into the database.
func (s Store) Create(ctx context.Context, task Task) error {
	const q = `
	INSERT INTO tasks
		(task_id, name, pid, wid, uid, estimated_seconds, active, date_created, date_updated, tracked_seconds)
	VALUES
		(:task_id, :name, :pid, :wid, :uid, :estimated_seconds, :active, :date_created, :date_updated, :tracked_seconds)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, task); err != nil {
		return fmt.Errorf("inserting task: %w", err)
	}

	return nil
}

// Update replaces a task document in the database.
func (s Store) Update(ctx context.Context, task Task) error {
	const q = `
	UPDATE
		tasks
	SET 
		"name" = :name,
		"estimated_seconds" = :estimated_seconds,
		"active" = :active,
		"date_updated" = :date_updated,
		"tracked_seconds" = :tracked_seconds
	WHERE
		task_id = :task_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, task); err != nil {
		return fmt.Errorf("updating taskID[%s]: %w", task.ID, err)
	}

	return nil
}

// Delete removes a task from the database.
func (s Store) Delete(ctx context.Context, taskID string) error {
	data := struct {
		TaskID string `db:"task_id"`
	}{
		TaskID: taskID,
	}

	const q = `
	DELETE FROM
		tasks
	WHERE
		task_id = :task_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting taskID[%s]: %w", taskID, err)
	}

	return nil
}

// Query retrieves a list of existing tasks from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Task, error) {
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
		tasks
	ORDER BY
		task_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tsks []Task
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tsks); err != nil {
		return nil, fmt.Errorf("selecting tasks: %w", err)
	}

	return tsks, nil
}

// QueryByID gets the specified task from the database.
func (s Store) QueryByID(ctx context.Context, taskID string) (Task, error) {
	data := struct {
		TaskID string `db:"task_id"`
	}{
		TaskID: taskID,
	}

	const q = `
	SELECT
		*
	FROM
		tasks
	WHERE 
		task_id = :task_id`

	var task Task
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &task); err != nil {
		return Task{}, fmt.Errorf("selecting taskID[%q]: %w", taskID, err)
	}

	return task, nil
}

// QueryUnique gets the specified project from the database.
func (s Store) QueryUnique(ctx context.Context, name, column, id string) string {
	data := struct {
		Name   string `db:"name"`
		Column string `db:"column"`
		Id     string `db:"id"`
	}{
		Name:   name,
		Column: column,
		Id:     id,
	}

	const q = `
	SELECT
		name
	FROM
		tasks
	WHERE 
		:column = :id AND name = :name`

	var nam string
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &nam); err != nil {
		return ""
	}

	return nam
}

// QueryBulkIDs gets all Tasks from the database.
func (s Store) QueryBulkIDs(ctx context.Context, taskID []string) ([]Task, error) {
	data := struct {
		TaskID []string `db:"task_id"`
	}{
		TaskID: taskID,
	}

	const q = `
	SELECT
		*
	FROM
		tasks
	WHERE 
		task_id in :task_id`

	var task []Task
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &task); err != nil {
		return nil, fmt.Errorf("selecting time_entry: %w", err)
	}

	return task, nil
}

//QueryProjectTasks retrieves a list of existing projects from the database.
func (s Store) QueryProjectTasks(ctx context.Context, projectID string, pageNumber, rowsPerPage int) ([]Task, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		ProjectID   string `db:"project_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		ProjectID:   projectID,
	}

	const q = `
	SELECT
		*
	FROM
		tasks
	WHERE 
		pid = :project_id
	ORDER BY
		name
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tasks []Task
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tasks); err != nil {
		return nil, fmt.Errorf("selecting projects: %w", err)
	}

	return tasks, nil
}

// QueryWorkspaceTasks retrieves a list of existing task from the database.
func (s Store) QueryWorkspaceTasks(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Task, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		WorkspaceID string `db:"workspace_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		WorkspaceID: workspaceID,
	}

	const q = `
	SELECT
		*
	FROM
		tasks
	WHERE
		wid = :workspace_id
	ORDER BY
		task_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tasks []Task
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tasks); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return tasks, nil
}

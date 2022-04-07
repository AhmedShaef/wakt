// Package db contains project related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

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

// Create inserts a new project into the database.
func (s Store) Create(ctx context.Context, project Project) error {
	const q = `
	INSERT INTO projects
		(project_id, name, wid, cid, active, is_private, template, template_id, billable, auto_estimates, estimated_hours, date_updated, color, rate, date_created)
	VALUES
		(:project_id, :name, :wid, :cid, :active, :is_private, :template, :template_id, :billable, :auto_estimates, :estimated_hours, :date_updated, :color, :rate, :date_created)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, project); err != nil {
		return fmt.Errorf("inserting project: %w", err)
	}

	return nil
}

// Update replaces a project document in the database.
func (s Store) Update(ctx context.Context, project Project) error {
	const q = `
	UPDATE
		projects
	SET 
		"name" = :name,
		"active" = :active,
		"is_private" = :is_private,
		"template" = :template,
		"template_id" = :template_id,
		"billable" = :billable,
		"auto_estimates" = :auto_estimates,
		"estimated_hours" = :estimated_hours,
		"date_updated" = :date_updated,
		"color" = :color,
		"rate" = :rate,
		"date_updated" = :date_updated
	WHERE
		project_id = :project_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, project); err != nil {
		return fmt.Errorf("updating projectID[%s]: %w", project.ID, err)
	}

	return nil
}

// Delete removes a project from the database.
func (s Store) Delete(ctx context.Context, projectID string) error {
	data := struct {
		projectID string `db:"project_id"`
	}{
		projectID: projectID,
	}

	const q = `
	DELETE FROM
		projects
	WHERE
		project_id = :project_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting projectID[%s]: %w", projectID, err)
	}

	return nil
}

//Query retrieves a list of existing projects from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Project, error) {
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
		projects
	ORDER BY
		project_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var projcts []Project
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &projcts); err != nil {
		return nil, fmt.Errorf("selecting projects: %w", err)
	}

	return projcts, nil
}

// QueryByID gets the specified project from the database.
func (s Store) QueryByID(ctx context.Context, projectID string) (Project, error) {
	data := struct {
		projectID string `db:"project_id"`
	}{
		projectID: projectID,
	}

	const q = `
	SELECT
		*
	FROM
		projects
	WHERE 
		project_id = :project_id`

	var projct Project
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &projct); err != nil {
		return Project{}, fmt.Errorf("selecting projectID[%q]: %w", projectID, err)
	}

	return projct, nil
}

// QueryUnique gets the specified project from the database.
func (s Store) QueryUnique(ctx context.Context, name, column, id string) string {
	data := struct {
		name   string `db:"name"`
		column string `db:"column"`
		id     string `db:"id"`
	}{
		name:   name,
		column: column,
		id:     id,
	}

	const q = `
	SELECT
		name
	FROM
		projects
	WHERE 
		:column = :id AND name = :name`

	var nam string
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &nam); err != nil {
		return ""
	}

	return nam
}

// QueryTrackedTime gets the specified project from the database.
func (s Store) QueryTrackedTime(ctx context.Context, projectID string) (time.Duration, error) {
	data := struct {
		projectID string `db:"project_id"`
	}{
		projectID: projectID,
	}

	const q = `
	SELECT
		SUM(tracked_seconds) AS tracked_seconds
	FROM
		tasks
	WHERE 
		pid = :project_id`

	var trackedSeconds time.Duration
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &trackedSeconds); err != nil {
		return 0, fmt.Errorf("selecting projectID[%q]: %w", projectID, err)
	}

	return trackedSeconds, nil
}

// QueryBulkIDs gets all Tasks from the database.
func (s Store) QueryBulkIDs(ctx context.Context, projectID []string) ([]Project, error) {
	data := struct {
		projectID []string `db:"project_id"`
	}{
		projectID: projectID,
	}

	const q = `
	SELECT
		*
	FROM
		projects
	WHERE 
		project_id = :project_id`

	var project []Project
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &project); err != nil {
		return nil, fmt.Errorf("selecting time_entrie: %w", err)
	}

	return project, nil
}

// QueryClientProjects retrieves a list of existing projects from the database.
func (s Store) QueryClientProjects(ctx context.Context, clientID string, pageNumber, rowsPerPage int) ([]Project, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		ClientID    string `db:"client_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		ClientID:    clientID,
	}

	const q = `
	SELECT
		*
	FROM
		projects
	WHERE 
		client_id = :client_id
	ORDER BY
		name
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var projcts []Project
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &projcts); err != nil {
		return nil, fmt.Errorf("selecting projects: %w", err)
	}

	return projcts, nil
}

// QueryWorkspaceProjects retrieves a list of existing project from the database.
func (s Store) QueryWorkspaceProjects(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Project, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		WorkspaceID string `db:"wid"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		WorkspaceID: workspaceID,
	}

	const q = `
	SELECT
		*
	FROM
		projects
	WHERE
		wid = :wid
	ORDER BY
		project_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var projcts []Project
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &projcts); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return projcts, nil
}

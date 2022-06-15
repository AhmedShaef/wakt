// Package db contains project user related CRUD functionality.
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

// Create inserts a new projectUser into the database.
func (s Store) Create(ctx context.Context, projectUser ProjectUser) error {
	const q = `
	INSERT INTO project_users
	   (project_user_id, pid, uid, wid, manager, date_created, date_updated)
	VALUES
	   (:project_user_id, :pid, :uid, :wid, :manager, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, projectUser); err != nil {
		return fmt.Errorf("inserting project user: %w", err)
	}

	return nil
}

// Update replaces a project user document in the database.
func (s Store) Update(ctx context.Context, projectUser ProjectUser) error {
	const q = `
	UPDATE
		project_users
	SET
		"manager" = :manager,
		"date_updated" = :date_updated
	WHERE
		"project_user_id" = :project_user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, projectUser); err != nil {
		return fmt.Errorf("updating projectUserID[%s]: %w", projectUser.ID, err)
	}

	return nil
}

// Delete removes a project user from the database.
func (s Store) Delete(ctx context.Context, projectUserID string) error {
	data := struct {
		ProjectUserID string `db:"project_user_id"`
	}{
		ProjectUserID: projectUserID,
	}

	const q = `
	DELETE FROM
		project_users
	WHERE
		project_user_id = :project_user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting projectUserID[%s]: %w", projectUserID, err)
	}

	return nil
}

// QueryWorkspaceProjectUsers retrieves a list of existing project user from the database.
func (s Store) QueryWorkspaceProjectUsers(ctx context.Context, WorkspaceID string, pageNumber, rowsPerPage int) ([]ProjectUser, error) {
	data := struct {
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
		WorkspaceID string `db:"workspace_id"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
		WorkspaceID: WorkspaceID,
	}

	const q = `
	SELECT
		*
	FROM
		project_users
	WHERE
		wid = :workspace_id
	ORDER BY
		project_user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var projectUsers []ProjectUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &projectUsers); err != nil {
		return nil, fmt.Errorf("selecting project user: %w", err)
	}

	return projectUsers, nil
}

// QueryByID gets the specified project user from the database.
func (s Store) QueryByID(ctx context.Context, projectUserID string) (ProjectUser, error) {
	data := struct {
		ProjectUserID string `db:"project_user_id"`
	}{
		ProjectUserID: projectUserID,
	}

	const q = `
	SELECT
		*
	FROM
		project_users
	WHERE
		project_user_id = :project_user_id`

	var projectUser ProjectUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &projectUser); err != nil {
		return ProjectUser{}, fmt.Errorf("selecting projectUserID[%q]: %w", projectUserID, err)
	}

	return projectUser, nil
}

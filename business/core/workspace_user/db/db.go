// Package db contains client related CRUD functionality.
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

// Invite creates new invitation.
func (s Store) Invite(ctx context.Context, workspaceUser WorkspaceUser) error {
	const q = `
	INSERT INTO workspace_users
		(workspace_id, uid, wid, active, invite_url)
	VALUES
		(:workspace_id, :uid, :wid, :active, :invite_url)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, workspaceUser); err != nil {
		return fmt.Errorf("inserting workspace user: %w", err)
	}

	return nil
}

// Update replaces a workspace user document in the database.
func (s Store) Update(ctx context.Context, workspaceUser WorkspaceUser) error {
	const q = `
	UPDATE
		workspace_users
	SET
		"active" = :active,
	WHERE
		workspace_user_id = :workspace_user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, workspaceUser); err != nil {
		return fmt.Errorf("updating workspaceUserID[%s]: %w", workspaceUser.ID, err)
	}

	return nil
}

// Delete removes a workspace user from the database.
func (s Store) Delete(ctx context.Context, workspaceUserID string) error {
	data := struct {
		WorkspaceUserID string `db:"workspace_user_id"`
	}{
		WorkspaceUserID: workspaceUserID,
	}

	const q = `
	DELETE FROM
		workspace_users
	WHERE
		workspace_user_id = :workspace_user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting workspaceUserID[%s]: %w", workspaceUserID, err)
	}

	return nil
}

// Query retrieves a list of existing workspace users from the database.
func (s Store) Query(ctx context.Context, userID string, pageNumber int, rowsPerPage int) ([]WorkspaceUser, error) {
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
		workspace_users
	WHERE
		uid = :user_id
	ORDER BY
		workspace_user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var workspaceUsers []WorkspaceUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &workspaceUsers); err != nil {
		return nil, fmt.Errorf("selecting workspace users: %w", err)
	}

	return workspaceUsers, nil
}

// QueryByID gets the specified workspace user from the database.
func (s Store) QueryByID(ctx context.Context, workspaceUserID string) (WorkspaceUser, error) {
	data := struct {
		WorkspaceUserID string `db:"workspace_user_id"`
	}{
		WorkspaceUserID: workspaceUserID,
	}

	const q = `
	SELECT
		*
	FROM
		workspace_users
	WHERE 
		workspace_user_id = :workspace_user_id`

	var workspaceUser WorkspaceUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &workspaceUser); err != nil {
		return WorkspaceUser{}, fmt.Errorf("selecting workspaceUserID[%q]: %w", workspaceUser, err)
	}

	return workspaceUser, nil
}

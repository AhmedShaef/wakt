// Package db contains group related CRUD functionality.
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

// Create inserts a new group into the database.
func (s Store) Create(ctx context.Context, group Group) error {
	const q = `
	INSERT INTO group
		(group_id, name, wid, date_updated)
	VALUES
		(:group_id, :name, :wid, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, group); err != nil {
		return fmt.Errorf("inserting group: %w", err)
	}

	return nil
}

// Update replaces a group document in the database.
func (s Store) Update(ctx context.Context, group Group) error {
	const q = `
	UPDATE
		group
	SET 
		"name" = :name,
		"date_updated" = :date_updated
	WHERE
		group_id = :group_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, group); err != nil {
		return fmt.Errorf("updating groupID[%s]: %w", group.ID, err)
	}

	return nil
}

// Delete removes a group from the database.
func (s Store) Delete(ctx context.Context, groupID string) error {
	data := struct {
		groupID string `db:"group_id"`
	}{
		groupID: groupID,
	}

	const q = `
	DELETE FROM
		group
	WHERE
		group_id = :group_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting groupID[%s]: %w", groupID, err)
	}

	return nil
}

// QueryByID gets the specified group from the database.
func (s Store) QueryByID(ctx context.Context, groupID string) (Group, error) {
	data := struct {
		groupID string `db:"group_id"`
	}{
		groupID: groupID,
	}

	const q = `
	SELECT
		*
	FROM
		group
	WHERE 
		group_id = :group_id`

	var group Group
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &group); err != nil {
		return Group{}, fmt.Errorf("selecting groupID[%q]: %w", groupID, err)
	}

	return group, nil
}

// QueryWorkspaceGroups retrieves a list of existing group from the database.
func (s Store) QueryWorkspaceGroups(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Group, error) {
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
		groups
	WHERE
		wid = :wid
	ORDER BY
		group_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var groups []Group
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &groups); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return groups, nil
}

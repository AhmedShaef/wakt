// Package db contains tag related CRUD functionality.
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

// Create inserts a new tag into the database.
func (s Store) Create(ctx context.Context, tag Tag) error {
	const q = `
	INSERT INTO tags
		(tag_id, name, wid, date_created, date_updated)
	VALUES
		(:tag_id, :name, :wid, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, tag); err != nil {
		return fmt.Errorf("inserting tag: %w", err)
	}

	return nil
}

// Update replaces a tag document in the database.
func (s Store) Update(ctx context.Context, tag Tag) error {
	const q = `
	UPDATE
		tags
	SET 
		"name" = :name,
		"date_updated" = :date_updated
	WHERE
		tag_id = :tag_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, tag); err != nil {
		return fmt.Errorf("updating tagID[%s]: %w", tag.ID, err)
	}

	return nil
}

// Delete removes a tag from the database.
func (s Store) Delete(ctx context.Context, tagID string) error {
	data := struct {
		TagID string `db:"tag_id"`
	}{
		TagID: tagID,
	}

	const q = `
	DELETE FROM
		tags
	WHERE
		tag_id = :tag_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting tagID[%s]: %w", tagID, err)
	}

	return nil
}

// QueryByID gets the specified tag from the database.
func (s Store) QueryByID(ctx context.Context, tagID string) (Tag, error) {
	data := struct {
		TagID string `db:"tag_id"`
	}{
		TagID: tagID,
	}

	const q = `
	SELECT
		*
	FROM
		tags
	WHERE 
		tag_id = :tag_id`

	var tag Tag
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &tag); err != nil {
		return Tag{}, fmt.Errorf("selecting tagID[%q]: %w", tagID, err)
	}

	return tag, nil
}

// QueryWorkspaceTags retrieves a list of existing tag from the database.
func (s Store) QueryWorkspaceTags(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Tag, error) {
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
		tags
	WHERE
		wid = :workspace_id
	ORDER BY
		tag_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var tags []Tag
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &tags); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return tags, nil
}

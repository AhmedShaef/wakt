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

// Create inserts a new client into the database.
func (s Store) Create(ctx context.Context, client Client) error {
	const q = `
	INSERT INTO clients
		(client_id, name, wid, notes, date_updated)
	VALUES
		(:client_id, :name, :wid, :notes, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, client); err != nil {
		return fmt.Errorf("inserting client: %w", err)
	}

	return nil
}

// Update replaces a client document in the database.
func (s Store) Update(ctx context.Context, client Client) error {
	const q = `
	UPDATE
		clients
	SET 
		"name" = :name,
		"notes" = :notes,
		"date_updated" = :date_updated
	WHERE
		client_id = :client_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, client); err != nil {
		return fmt.Errorf("updating clientID[%s]: %w", client.ID, err)
	}

	return nil
}

// Delete removes a client from the database.
func (s Store) Delete(ctx context.Context, clientID string) error {
	data := struct {
		clientID string `db:"client_id"`
	}{
		clientID: clientID,
	}

	const q = `
	DELETE FROM
		clients
	WHERE
		client_id = :client_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting clientID[%s]: %w", clientID, err)
	}

	return nil
}

// Query retrieves a list of existing client from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Client, error) {
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
		clients
	ORDER BY
		client_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var clients []Client
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &clients); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return clients, nil
}

// QueryByID gets the specified client from the database.
func (s Store) QueryByID(ctx context.Context, clientID string) (Client, error) {
	data := struct {
		clientID string `db:"client_id"`
	}{
		clientID: clientID,
	}

	const q = `
	SELECT
		*
	FROM
		clients
	WHERE 
		client_id = :client_id`

	var client Client
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &client); err != nil {
		return Client{}, fmt.Errorf("selecting clientID[%q]: %w", clientID, err)
	}

	return client, nil
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
		clients
	WHERE 
		:column = :id AND name = :name`

	var nam string
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &nam); err != nil {
		return ""
	}

	return nam
}

// QueryWorkspaceClients retrieves a list of existing client from the database.
func (s Store) QueryWorkspaceClients(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Client, error) {
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
		clients
	WHERE
		wid = :wid
	ORDER BY
		client_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var clients []Client
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &clients); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return clients, nil
}

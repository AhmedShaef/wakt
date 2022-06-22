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

// Create inserts a new team into the database.
func (s Store) Create(ctx context.Context, team Team) error {
	const q = `
	INSERT INTO teams
	   (team_id, pid, uid, wid, manager, date_created, date_updated)
	VALUES
	   (:team_id, :pid, :uid, :wid, :manager, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, team); err != nil {
		return fmt.Errorf("inserting project user: %w", err)
	}

	return nil
}

// Update replaces a project user document in the database.
func (s Store) Update(ctx context.Context, team Team) error {
	const q = `
	UPDATE
		teams
	SET
		"manager" = :manager,
		"date_updated" = :date_updated
	WHERE
		"team_id" = :team_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, team); err != nil {
		return fmt.Errorf("updating teamID[%s]: %w", team.ID, err)
	}

	return nil
}

// Delete removes a project user from the database.
func (s Store) Delete(ctx context.Context, teamID string) error {
	data := struct {
		TeamID string `db:"team_id"`
	}{
		TeamID: teamID,
	}

	const q = `
	DELETE FROM
		teams
	WHERE
		team_id = :team_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting teamID[%s]: %w", teamID, err)
	}

	return nil
}

// QueryWorkspaceTeams retrieves a list of existing project user from the database.
func (s Store) QueryWorkspaceTeams(ctx context.Context, WorkspaceID string, pageNumber, rowsPerPage int) ([]Team, error) {
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
		teams
	WHERE
		wid = :workspace_id
	ORDER BY
		team_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var teams []Team
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &teams); err != nil {
		return nil, fmt.Errorf("selecting project user: %w", err)
	}

	return teams, nil
}

// QueryByID gets the specified project user from the database.
func (s Store) QueryByID(ctx context.Context, teamID string) (Team, error) {
	data := struct {
		TeamID string `db:"team_id"`
	}{
		TeamID: teamID,
	}

	const q = `
	SELECT
		*
	FROM
		teams
	WHERE
		team_id = :team_id`

	var team Team
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &team); err != nil {
		return Team{}, fmt.Errorf("selecting teamID[%q]: %w", teamID, err)
	}

	return team, nil
}

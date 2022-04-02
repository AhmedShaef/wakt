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

// Create inserts a new workspace into the database.
func (s Store) Create(ctx context.Context, workspace Workspace) error {
	const q = `
	INSERT INTO workspaces
		(workspace_id, name, profile, premium, admin, default_hourly_rate, default_currency, only_admin_may_create_projects, only_admin_see_billable_rates, only_admin_see_team_dashboard, project_billable_by_default, rounding, rounding_minutes, date_updated, logo_url, ical_url, ical_enabled)
	VALUES
		(:workspace_id, :name, :profile, :premium, :admin, :default_hourly_rate, :default_currency, :only_admin_may_create_projects, :only_admin_see_billable_rates, :only_admin_see_team_dashboard, :project_billable_by_default, :rounding, :rounding_minutes, :date_updated, :logo_url, :ical_url, :ical_enabled)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, workspace); err != nil {
		return fmt.Errorf("inserting workspace: %w", err)
	}

	return nil
}

// Update replaces a workspace document in the database.
func (s Store) Update(ctx context.Context, workspace Workspace) error {
	const q = `
	UPDATE
		users
	SET
		name = :name,
		profile = :profile,
		premium = :premium,
		admin = :admin,
		default_hourly_rate = :default_hourly_rate,
		default_currency = :default_currency,
		only_admin_may_create_projects = :only_admin_may_create_projects,
		only_admin_see_billable_rates = :only_admin_see_billable_rates,
		only_admin_see_team_dashboard = :only_admin_see_team_dashboard,
		project_billable_by_default = :project_billable_by_default,
		rounding = :rounding,
		rounding_minutes = :rounding_minutes,
		date_updated = :date_updated,
		logo_url = :logo_url,
		ical_url = :ical_url,
		ical_enabled = :ical_enabled
	WHERE
		workspace_id = :workspace_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, workspace); err != nil {
		return fmt.Errorf("updating workspaceID[%s]: %w", workspace.ID, err)
	}

	return nil

}

// Delete removes a workspace from the database.
func (s Store) Delete(ctx context.Context, workspaceID string) error {
	data := struct {
		workspaceID string `db:"workspace_id"`
	}{
		workspaceID: workspaceID,
	}

	const q = `
	DELETE FROM
		workspaces
	WHERE
		workspace_id = :workspace_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting workspaceID[%s]: %w", workspaceID, err)
	}

	return nil
}

// Query retrieves a list of existing workspaces from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Workspace, error) {
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
		workspaces
	ORDER BY
		workspace_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var workspaces []Workspace
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &workspaces); err != nil {
		return nil, fmt.Errorf("selecting workspaces: %w", err)
	}

	return workspaces, nil
}

// QueryByID gets the specified workspace from the database.
func (s Store) QueryByID(ctx context.Context, workspaceID string) (Workspace, error) {
	data := struct {
		workspaceID string `db:"workspace_id"`
	}{
		workspaceID: workspaceID,
	}

	const q = `
	SELECT
		*
	FROM
		workspaces
	WHERE 
		workspace_id = :workspace_id`

	var workspace Workspace
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &workspace); err != nil {
		return Workspace{}, fmt.Errorf("selecting userID[%q]: %w", workspaceID, err)
	}

	return workspace, nil
}

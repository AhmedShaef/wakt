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

// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, user User) error {
	const q = `
	INSERT INTO clients
		(user_id, api_token, default_wid, email, password_hash, full_name, jquery_time_of_day_format, jquery_date_format, time_of_day_format, date_format, store_start_and_stop_time, beginning_of_week, language, image_url, sidebar_piechart, date_created, date_updated, record_timeline,should_upgrade, new_blog_post, send_product_emails, send_weekly_report, send_timer_notifications, openid_enabled, timezone, invitation, duration_format)
	VALUES
		(:user_id, :api_token, :default_wid, :email, :password_hash, :full_name, :jquery_time_of_day_format, :jquery_date_format, :time_of_day_format, :date_format, :store_start_and_stop_time, :beginning_of_week, :language, :image_url, :sidebar_piechart, :date_created, :date_updated, :record_timeline,should_upgrade, :new_blog_post, :send_product_emails, :send_weekly_report, :send_timer_notifications, :openid_enabled, :timezone, :invitation, :duration_format)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, user); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (s Store) Update(ctx context.Context, user User) error {
	const q = `
	UPDATE
		users
	SET 
		api_token = :api_token,
		default_wid = :default_wid,
		email = :email,
		password_hash = :password_hash,
		full_name = :full_name,
		jquery_time_of_day_format = :jquery_time_of_day_format,
		jquery_date_format = :jquery_date_format,
		time_of_day_format = :time_of_day_format,
		date_format = :date_format,
		store_start_and_stop_time = :store_start_and_stop_time,
		beginning_of_week = :beginning_of_week,
		language = :language,
		image_url = :image_url,
		sidebar_piechart = :sidebar_piechart,
		date_created = :date_created,
		date_updated = :date_updated,
		record_timeline = :record_timeline,
		should_upgrade = :should_upgrade,
		new_blog_post = :new_blog_post,
		send_product_emails = :send_product_emails,
		send_weekly_report = :send_weekly_report,
		send_timer_notifications = :send_timer_notifications,
		openid_enabled = :openid_enabled,
		timezone = :timezone,
		invitation = :invitation,
		duration_format = :duration_format
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, user); err != nil {
		return fmt.Errorf("updating userID[%s]: %w", user.ID, err)
	}

	return nil
}

// Delete removes a user from the database.
func (s Store) Delete(ctx context.Context, userID string) error {
	data := struct {
		userID string `db:"user_id"`
	}{
		userID: userID,
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting userID[%s]: %w", userID, err)
	}

	return nil
}

// QueryByID gets the specified user from the database.
func (s Store) QueryByID(ctx context.Context, userID string) (User, error) {
	data := struct {
		userID string `db:"user_id"`
	}{
		userID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = :user_id`

	var user User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &user); err != nil {
		return User{}, fmt.Errorf("selecting userID[%q]: %w", userID, err)
	}

	return user, nil
}

// QueryByEmail gets the specified user from the database by email.
func (s Store) QueryByEmail(ctx context.Context, email string) (User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`

	var user User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &user); err != nil {
		return User{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	return user, nil
}

// QueryWorkspaceUsers retrieves a list of existing user from the database.
func (s Store) QueryWorkspaceUsers(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]User, error) {
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
		users
	WHERE
		wid = :wid
	ORDER BY
		user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var users []User
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &users); err != nil {
		return nil, fmt.Errorf("selecting client: %w", err)
	}

	return users, nil
}

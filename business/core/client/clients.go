// Package client provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/client/db"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrInvalidEmail          = errors.New("email is not valid")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrInvalidPassword       = errors.New("password is not valid")
)

// Core manages the set of APIs for user access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create inserts a new client into the database.
func (c Core) Create(ctx context.Context, nc NewClient, email string, now time.Time) (Client, error) {
	if err := validate.Check(nc); err != nil {
		return Client{}, fmt.Errorf("validating data: %w", err)
	}

	usr, err := user.Core{}.QueryByEmail(ctx, email)
	if err != nil {
		return Client{}, fmt.Errorf("querying user: %w", err)
	}

	if nc.Wid == "" {
		nc.Wid = usr.DefaultWid
	}

	nameInWorkspace := c.store.QueryUnique(ctx, nc.Name, "wid", nc.Wid)
	if nameInWorkspace != "" {
		return Client{}, fmt.Errorf("project name is not unique for workspace")
	}

	dbclint := db.Client{
		ID:          validate.GenerateID(),
		Name:        nc.Name,
		Wid:         nc.Wid,
		Notes:       nc.Notes,
		DateUpdated: now,
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		if err := c.store.Tran(tx).Create(ctx, dbclint); err != nil {
			if errors.Is(err, database.ErrDBDuplicatedEntry) {
				return fmt.Errorf("create: %w", ErrUniqueEmail)
			}
			return fmt.Errorf("create: %w", err)
		}
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return Client{}, fmt.Errorf("tran: %w", err)
	}

	return toClient(dbclint), nil
}

// Update replaces a client document in the database.
func (c Core) Update(ctx context.Context, clientID string, uc UpdateClient, now time.Time) error {
	if err := validate.CheckID(clientID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uc); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbclient, err := c.store.QueryByID(ctx, clientID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating client clientID[%s]: %w", clientID, err)
	}

	if uc.Name != nil {
		dbclient.Name = *uc.Name
	}
	if uc.Notes != nil {
		dbclient.Notes = *uc.Notes
	}
	dbclient.DateUpdated = now

	if err := c.store.Update(ctx, dbclient); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a client from the database.
func (c Core) Delete(ctx context.Context, clientID string) error {
	if err := validate.CheckID(clientID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, clientID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

//Query retrieves a list of existing client from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Client, error) {
	dbclient, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toClientsSlice(dbclient), nil
}

// QueryByID gets the specified client from the database.
func (c Core) QueryByID(ctx context.Context, clientID string) (Client, error) {
	if err := validate.CheckID(clientID); err != nil {
		return Client{}, ErrInvalidID
	}

	dbclient, err := c.store.QueryByID(ctx, clientID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Client{}, ErrNotFound
		}
		return Client{}, fmt.Errorf("query: %w", err)
	}

	return toClient(dbclient), nil
}
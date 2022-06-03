// Package projects provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.

package project

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/project/db"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("user not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
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

// Create inserts a new project into the database.
func (c Core) Create(ctx context.Context, userID string, np NewProject, now time.Time) (Project, error) {
	if err := validate.Check(np); err != nil {
		return Project{}, fmt.Errorf("validating data: %w", err)
	}

	nameInWorkspace := c.store.QueryUnique(ctx, np.Name, "wid", np.Wid)
	if nameInWorkspace != "" {
		return Project{}, fmt.Errorf("project name is not unique for workspace")
	}

	nameInClient := c.store.QueryUnique(ctx, np.Name, "cid", np.Cid)
	if nameInClient != "" {
		return Project{}, fmt.Errorf("project name is not unique for client")
	}

	// Set defaults
	dbprojct := db.Project{
		IsPrivate:      true,
		Billable:       true,
		AutoEstimates:  false,
		EstimatedHours: 0,
		Rate:           0.0,
		HexColor:       "#ffffff",
	}
	// Set values from NewProject
	dbprojct = db.Project{
		ID:             validate.GenerateID(),
		Name:           np.Name,
		Wid:            np.Wid,
		Cid:            np.Cid,
		Uid:            userID,
		Active:         false,
		IsPrivate:      np.IsPrivate,
		Billable:       np.Billable,
		AutoEstimates:  np.AutoEstimates,
		EstimatedHours: np.EstimatedHours,
		DateCreated:    now,
		DateUpdated:    now,
		Rate:           np.Rate,
		HexColor:       np.HexColor,
	}

	if err := c.store.Create(ctx, dbprojct); err != nil {
		return Project{}, fmt.Errorf("create: %w", err)
	}

	return toProject(dbprojct), nil
}

// Update replaces a project document in the database.
func (c Core) Update(ctx context.Context, projectID string, up UpdateProject, now time.Time) error {
	if err := validate.CheckID(projectID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbprojct, err := c.store.QueryByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating project projectID[%s]: %w", projectID, err)
	}

	if up.Name != nil {
		dbprojct.Name = *up.Name
	}
	if up.Active != nil {
		dbprojct.Active = *up.Active
	}
	if up.IsPrivate != nil {
		dbprojct.IsPrivate = *up.IsPrivate
	}
	if up.AutoEstimates != nil {
		dbprojct.AutoEstimates = *up.AutoEstimates
	}
	if !dbprojct.AutoEstimates {
		if up.EstimatedHours != nil {
			dbprojct.EstimatedHours = *up.EstimatedHours
		}
	}
	if up.Billable != nil {
		dbprojct.Billable = *up.Billable
	}
	if up.Rate != nil {
		dbprojct.Rate = *up.Rate
	}
	if up.HexColor != nil {
		dbprojct.HexColor = *up.HexColor
	}
	dbprojct.DateUpdated = now

	if err := c.store.Update(ctx, dbprojct); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a project from the database.
func (c Core) Delete(ctx context.Context, projectID string) error {
	if err := validate.CheckID(projectID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, projectID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByID gets the specified project from the database.
func (c Core) QueryByID(ctx context.Context, projectID string) (Project, error) {
	if err := validate.CheckID(projectID); err != nil {
		return Project{}, ErrInvalidID
	}

	dbprojct, err := c.store.QueryByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Project{}, ErrNotFound
		}
		return Project{}, fmt.Errorf("query: %w", err)
	}

	return toProject(dbprojct), nil
}

// QueryClientProjects retrieves a list of existing projects from the database.
func (c Core) QueryClientProjects(ctx context.Context, clientID string, pageNumber, rowsPerPage int) ([]Project, error) {
	dbprojects, err := c.store.QueryClientProjects(ctx, clientID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toProjectsSlice(dbprojects), nil
}

// QueryWorkspaceProjects retrieves a list of existing workspace from the database.
func (c Core) QueryWorkspaceProjects(ctx context.Context, workspaceID string, pageNumber, rowsPerPage int) ([]Project, error) {
	dbProjects, err := c.store.QueryWorkspaceProjects(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return toProjectsSlice(dbProjects), nil
}

// QueryUserProjects retrieves a list of existing projects from the database.
func (c Core) QueryUserProjects(ctx context.Context, userID string, pageNumber, rowsPerPage int) ([]Project, error) {
	if err := validate.CheckID(userID); err != nil {
		return []Project{}, ErrInvalidID
	}
	dbprojects, err := c.store.QueryUserProjects(ctx, userID, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toProjectsSlice(dbprojects), nil
}

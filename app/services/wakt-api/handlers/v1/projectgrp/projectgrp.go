// Package projectgrp maintains the group of handlers for user access.
package projectgrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/client"
	"github.com/AhmedShaef/wakt/business/core/project"
	"github.com/AhmedShaef/wakt/business/core/task"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
	"strconv"
	"strings"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	Project   project.Core
	Workspace workspace.Core
	User      user.Core
	Task      task.Core
}

// Create adds a new user to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var np project.NewProject
	if err := web.Decode(r, &np); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaces, err := h.Workspace.QueryByID(ctx, np.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", np.Wid, err)
		}
	}

	// If you are not an admin and looking to update a client you don't own.
	if workspaces.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}
	prj, err := h.Project.Create(ctx, np, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("user[%+v]: %w", &prj, err)
		}
	}

	return web.Respond(ctx, w, prj, http.StatusCreated)
}

// Update updates a user in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var up project.UpdateProject
	if err := web.Decode(r, &up); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	projectID := web.Param(r, "id")

	projects, err := h.Project.QueryByID(ctx, projectID)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", projectID, err)
		}
	}
	workspaces, err := h.Workspace.QueryByID(ctx, projects.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaces.ID, err)
		}
	}
	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Project.Update(ctx, projectID, up, v.Now); err != nil {
		switch {
		case errors.Is(err, project.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] User[%+v]: %w", projectID, &up, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// DeleteByID removes a user from the system.
func (h Handlers) DeleteByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	projectID := web.Param(r, "id")

	projects, err := h.Project.QueryByID(ctx, projectID)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", projectID, err)
		}
	}

	workspaces, err := h.Workspace.QueryByID(ctx, projects.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaces.ID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Project.Delete(ctx, projectID); err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", projectID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// QueryByID returns a user by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	projectID := web.Param(r, "id")

	projects, err := h.Project.QueryByID(ctx, projectID)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, client.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", projectID, err)
		}
	}

	workspaces, err := h.Workspace.QueryByID(ctx, projects.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaces.ID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	prj, err := h.Project.QueryByID(ctx, projectID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", projectID, err)
		}
	}

	return web.Respond(ctx, w, prj, http.StatusOK)
}

// Delete removes a user from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	projctID := web.Param(r, "id")
	projectIDs := strings.Split(projctID, ",")

	for _, projectID := range projectIDs {
		projects, err := h.Project.QueryByID(ctx, projectID)
		if err != nil {
			switch {
			case errors.Is(err, project.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, project.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", projectID, err)
			}
		}

		workspaces, err := h.Workspace.QueryByID(ctx, projects.Wid)
		if err != nil {
			switch {
			case errors.Is(err, workspace.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, workspace.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", workspaces.ID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if claims.Subject != workspaces.Uid {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.Project.Delete(ctx, projectID); err != nil {
			switch {
			case errors.Is(err, user.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			default:
				return fmt.Errorf("ID[%s]: %w", projectID, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// QueryProjectTasks returns a list of workspaces with paging.
func (h Handlers) QueryProjectTasks(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	taskID := web.Param(r, "id")
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format, page[%s]", page), http.StatusBadRequest)
	}
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format, rows[%s]", rows), http.StatusBadRequest)
	}

	tasks, err := h.Workspace.QueryByID(ctx, taskID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, task.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying task[%s]: %w", taskID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != tasks.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	tsk, err := h.Task.QueryProjectTasks(ctx, taskID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for task: %w", err)
	}

	return web.Respond(ctx, w, tsk, http.StatusOK)
}

// QueryProjectUsers returns a list of workspaces with paging.
func (h Handlers) QueryProjectUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	userID := web.Param(r, "id")
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format, page[%s]", page), http.StatusBadRequest)
	}
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format, rows[%s]", rows), http.StatusBadRequest)
	}

	users, err := h.Workspace.QueryByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying user[%s]: %w", userID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != users.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	usr, err := h.User.QueryProjectUsers(ctx, userID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for users: %w", err)
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}
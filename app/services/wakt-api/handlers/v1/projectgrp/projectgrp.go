// Package projectgrp maintains the group of handlers for project access.
package projectgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AhmedShaef/wakt/business/core/project"
	"github.com/AhmedShaef/wakt/business/core/task"
	"github.com/AhmedShaef/wakt/business/core/team"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/core/workspace_user"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
)

// Handlers manages the set of project endpoints.
type Handlers struct {
	Project       project.Core
	ProjectUser   team.Core
	Workspace     workspace.Core
	WorkspaceUser workspace_user.Core
	User          user.Core
	Task          task.Core
}

// Create adds a new project to the system.
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

	if np.Wid == "" {
		users, err := h.User.QueryByID(ctx, claims.Subject)
		if err != nil {
			return fmt.Errorf("unable to query user: %w", err)
		}
		np.Wid = users.DefaultWid
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

	workspaceUser, err := h.WorkspaceUser.QueryByuIDwID(ctx, np.Wid, claims.Subject)
	if err != nil {
		switch {
		case errors.Is(err, workspace_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace user[%s]: %w", workspaceUser.ID, err)
		}
	}

	// If you are not an admin and looking to update a project you don't own.
	if workspaces.OnlyAdminMayCreateProjects {
		if !workspaceUser.Admin || workspaces.Uid != claims.Subject {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}
	} else {
		if workspaces.Uid != claims.Subject {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}
	}

	prj, err := h.Project.Create(ctx, claims.Subject, np, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("project[%+v]: %w", &prj, err)
		}
	}
	npu := team.NewProjectUser{
		Pid:     prj.ID,
		Uid:     workspaces.Uid,
		Wid:     prj.Wid,
		Manager: true,
	}
	projectUser, err := h.ProjectUser.Create(ctx, npu, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, team.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("team[%+v]: %w", &projectUser, err)
		}
	}

	return web.Respond(ctx, w, prj, http.StatusCreated)
}

// Update updates a project in the system.
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

// QueryByID returns a project by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	workspaceUser, err := h.WorkspaceUser.QueryByuIDwID(ctx, workspaces.ID, workspaces.Uid)
	if err != nil {
		switch {
		case errors.Is(err, workspace_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace user[%s]: %w", workspaceUser.ID, err)
		}
	}

	// If you are not an admin and looking to update a project you don't own.
	if workspaces.OnlyAdminSeeBillableRates {
		if !workspaceUser.Admin {
			projects.Rate = 0.0
		}
	}
	if workspaces.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	return web.Respond(ctx, w, projects, http.StatusOK)
}

// BulkDelete removes a project from the system.
func (h Handlers) BulkDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
			case errors.Is(err, project.ErrInvalidID):
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

	projectID := web.Param(r, "id")
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

	projects, err := h.Project.QueryByID(ctx, projectID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, task.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying task[%s]: %w", projectID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != projects.UID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	tsk, err := h.Task.QueryProjectTasks(ctx, projectID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for task: %w", err)
	}

	return web.Respond(ctx, w, tsk, http.StatusOK)
}

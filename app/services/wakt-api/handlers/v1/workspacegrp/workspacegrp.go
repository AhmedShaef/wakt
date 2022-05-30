// Package workspacegrp maintains the group of handlers for workspace access.
package workspacegrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/app/tooling/uploader"
	"github.com/AhmedShaef/wakt/business/core/client"
	"github.com/AhmedShaef/wakt/business/core/group"
	"github.com/AhmedShaef/wakt/business/core/project"
	"github.com/AhmedShaef/wakt/business/core/project_user"
	"github.com/AhmedShaef/wakt/business/core/tag"
	"github.com/AhmedShaef/wakt/business/core/task"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/core/workspace_user"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
	"strconv"
)

// Handlers manages the set of workspace endpoints.
type Handlers struct {
	Workspace     workspace.Core
	User          user.Core
	Client        client.Core
	Group         group.Core
	Project       project.Core
	Task          task.Core
	Tag           tag.Core
	ProjectUser   project_user.Core
	WorkspaceUser workspace_user.Core
}

// Create adds a new workspace to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var nw workspace.NewWorkspace
	if err := web.Decode(r, &nw); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	// If you are not an admin and looking to update a workspace you don't own.
	if nw.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.Workspace.Create(ctx, nw, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("workspace[%+v]: %w", &work, err)
		}
	}

	workspaceUser, err := h.WorkspaceUser.Create(ctx, work.ID, work.Uid, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, workspace_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("workspace[%+v]: %w", &workspaceUser, err)
		}
	}
	return web.Respond(ctx, w, work, http.StatusCreated)
}

// Update updates a workspace in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd workspace.UpdateWorkspace
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaceID := web.Param(r, "id")

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Workspace.Update(ctx, workspaceID, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] Workspace[%+v]: %w", workspaceID, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// UpdateLogo updates a workspace logo in the system.
func (h Handlers) UpdateLogo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	if err := r.ParseMultipartForm(2 << 10); err != nil {
		return fmt.Errorf("unable to parse multipart form: %w", err)
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	file, handler, err := r.FormFile("profileLogo")
	if err != nil {
		return fmt.Errorf("unable to get file %S: %w", handler.Filename, err)
	}
	defer file.Close()

	name, err := uploader.UploadImage(file)
	if err != nil {
		return fmt.Errorf("unable to upload logo: %w", err)
	}
	upd := workspace.UpdateWorkspace{
		LogoURL: name,
	}

	if err := h.Workspace.UpdateLogo(ctx, workspaces.ID, upd, v.Now); err != nil {
		return fmt.Errorf("ID[%s] User[%+v]: %w", claims.Subject, &upd, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Query returns a list of workspaces with paging.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

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

	work, err := h.Workspace.Query(ctx, claims.Subject, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryByID returns a workspace by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	return web.Respond(ctx, w, workspaces, http.StatusOK)
}

// QueryWorkspaceUsers returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.WorkspaceUser.Query(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryWorkspaceClients returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceClients(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.Client.QueryWorkspaceClients(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryWorkspaceGroups returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceGroups(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.Group.QueryWorkspaceGroups(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryWorkspaceProjects returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceProjects(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.Project.QueryWorkspaceProjects(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryWorkspaceTasks returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceTasks(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.Task.QueryWorkspaceTasks(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryWorkspaceTags returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceTags(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.Tag.QueryWorkspaceTags(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// QueryWorkspaceProjectUsers returns a list of workspaces with paging.
func (h Handlers) QueryWorkspaceProjectUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceID := web.Param(r, "id")
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

	workspaces, err := h.Workspace.QueryByID(ctx, workspaceID)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != workspaces.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	work, err := h.ProjectUser.QueryWorkspaceProjectUsers(ctx, workspaceID, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for workspaces: %w", err)
	}

	return web.Respond(ctx, w, work, http.StatusOK)
}

// Package projectusergrp maintains the group of handlers for projectuser access.
package projectusergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/project_user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
	"strings"
)

// Handlers manages the set of projectuser endpoints.
type Handlers struct {
	ProjectUser project_user.Core
	Workspace   workspace.Core
}

// Create adds a new project_user to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var npu project_user.NewProjectUser
	if err := web.Decode(r, &npu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	// If you are not an admin and looking to update a project_user you don't own.
	if !claims.Authorized(auth.RoleAdmin) && npu.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	clint, err := h.ProjectUser.Create(ctx, npu, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, project_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("project_user[%+v]: %w", &clint, err)
		}
	}

	return web.Respond(ctx, w, clint, http.StatusCreated)
}

// UpdateByID updates a project_user in the system.
func (h Handlers) UpdateByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd project_user.UpdateProjectUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	projectUserID := web.Param(r, "id")

	projectUsers, err := h.ProjectUser.QueryByID(ctx, projectUserID)
	if err != nil {
		switch {
		case errors.Is(err, project_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", projectUserID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != projectUsers.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.ProjectUser.Update(ctx, projectUserID, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, project_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] ProjectUser[%+v]: %w", projectUserID, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// DeleteByID removes a project_user from the system.
func (h Handlers) DeleteByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	projectUserID := web.Param(r, "id")

	projectUsers, err := h.ProjectUser.QueryByID(ctx, projectUserID)
	if err != nil {
		switch {
		case errors.Is(err, project_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, project_user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", projectUserID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != projectUsers.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.ProjectUser.Delete(ctx, projectUserID); err != nil {
		switch {
		case errors.Is(err, project_user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", projectUserID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Update updates a project_user in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd project_user.UpdateProjectUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	projctUserID := web.Param(r, "id")
	projectUserIDs := strings.Split(projctUserID, ",")

	for _, projectUserID := range projectUserIDs {
		projectUsers, err := h.ProjectUser.QueryByID(ctx, projectUserID)
		if err != nil {
			switch {
			case errors.Is(err, project_user.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, project_user.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", projectUserID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if !claims.Authorized(auth.RoleAdmin) && claims.Subject != projectUsers.Uid {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.ProjectUser.Update(ctx, projectUserID, upd, v.Now); err != nil {
			switch {
			case errors.Is(err, project_user.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, project_user.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("ID[%s] ProjectUser[%+v]: %w", projectUserID, &upd, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a project_user from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	projctUserID := web.Param(r, "id")
	projectUserIDs := strings.Split(projctUserID, ",")

	for _, projectUserID := range projectUserIDs {
		projectUsers, err := h.ProjectUser.QueryByID(ctx, projectUserID)
		if err != nil {
			switch {
			case errors.Is(err, project_user.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, project_user.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", projectUserID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if !claims.Authorized(auth.RoleAdmin) && claims.Subject != projectUsers.Uid {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.ProjectUser.Delete(ctx, projectUserID); err != nil {
			switch {
			case errors.Is(err, project_user.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			default:
				return fmt.Errorf("ID[%s]: %w", projectUserID, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

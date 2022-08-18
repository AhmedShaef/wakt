// Package workspaceusergrp maintains the group of handlers for workspaceuser access.
package workspaceusergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/AhmedShaef/wakt/business/core/workspaceuser"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
)

// Handlers manages the set of workspaceUser endpoints.
type Handlers struct {
	WorkspaceUser workspaceuser.Core
}

// Invite adds a new workspaceUser to the system.
func (h Handlers) Invite(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var nwu workspaceuser.InviteUsers
	if err := web.Decode(r, &nwu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaceUser, err := h.WorkspaceUser.QueryByID(ctx, nwu.InviterID)
	if err != nil {
		switch {
		case errors.Is(err, workspaceuser.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspaceuser.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspaceUser[%s]: %w", workspaceUser.ID, err)
		}
	}

	// If you are not an admin and looking to update a workspaceUser you don't own.
	if !workspaceUser.Admin && workspaceUser.UID != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	clint, err := h.WorkspaceUser.InviteUser(ctx, workspaceUser.WID, nwu, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, workspaceuser.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspaceuser.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("workspaceUser[%+v]: %w", &clint, err)
		}
	}

	return web.Respond(ctx, w, clint, http.StatusCreated)
}

// Update updates a workspaceUser in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd workspaceuser.UpdateWorkspaceUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaceUserID := web.Param(r, "id")

	workspaceUsers, err := h.WorkspaceUser.QueryByID(ctx, workspaceUserID)
	if err != nil {
		switch {
		case errors.Is(err, workspaceuser.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspaceuser.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceUserID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !workspaceUsers.Admin && claims.Subject != workspaceUsers.UID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.WorkspaceUser.Update(ctx, workspaceUserID, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, workspaceuser.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspaceuser.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] WorkspaceUser[%+v]: %w", workspaceUserID, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a workspaceUser from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	workspaceUserID := web.Param(r, "id")

	workspaceUsers, err := h.WorkspaceUser.QueryByID(ctx, workspaceUserID)
	if err != nil {
		switch {
		case errors.Is(err, workspaceuser.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspaceuser.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", workspaceUserID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !workspaceUsers.Admin && claims.Subject != workspaceUsers.UID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.WorkspaceUser.Delete(ctx, workspaceUserID); err != nil {
		switch {
		case errors.Is(err, workspaceuser.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", workspaceUserID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

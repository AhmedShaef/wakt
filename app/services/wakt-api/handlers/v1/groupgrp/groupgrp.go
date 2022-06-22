// Package groupgrp maintains the group of handlers for group access.
package groupgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/AhmedShaef/wakt/business/core/group"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
)

// Handlers manages the set of group endpoints.
type Handlers struct {
	Group     group.Core
	Workspace workspace.Core
	User      user.Core
}

// Create adds a new group to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ng group.NewGroup
	if err := web.Decode(r, &ng); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if ng.Wid == "" {
		users, err := h.User.QueryByID(ctx, claims.Subject)
		if err != nil {
			return fmt.Errorf("unable to query user:%w", err)
		}
		ng.Wid = users.DefaultWid
	} else {
		workspaces, err := h.Workspace.QueryByID(ctx, ng.Wid)
		if err != nil {
			return fmt.Errorf("unable to query workspace:%w", err)
		}

		// If you are not an admin and looking to update a tag you don't own.
		if workspaces.Uid != claims.Subject {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}
	}

	groups, err := h.Group.Create(ctx, claims.Subject, ng, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, group.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, group.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("group[%+v]: %w", &groups, err)
		}
	}

	return web.Respond(ctx, w, groups, http.StatusCreated)
}

// Update updates a group in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ug group.UpdateGroup
	if err := web.Decode(r, &ug); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	groupID := web.Param(r, "id")

	groups, err := h.Group.QueryByID(ctx, groupID)
	if err != nil {
		switch {
		case errors.Is(err, group.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, group.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", groupID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != groups.UID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Group.Update(ctx, groupID, ug, v.Now); err != nil {
		switch {
		case errors.Is(err, group.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, group.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] group[%+v]: %w", groupID, &ug, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a group from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	groupID := web.Param(r, "id")

	groups, err := h.Group.QueryByID(ctx, groupID)
	if err != nil {
		switch {
		case errors.Is(err, group.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, group.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", groupID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != groups.UID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Group.Delete(ctx, groupID); err != nil {
		switch {
		case errors.Is(err, group.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", groupID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

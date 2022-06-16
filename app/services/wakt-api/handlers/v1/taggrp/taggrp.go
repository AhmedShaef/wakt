// Package taggrp maintains the tag of handlers for tag access.
package taggrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/tag"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
)

// Handlers manages the set of tag endpoints.
type Handlers struct {
	Tag       tag.Core
	Workspace workspace.Core
	User      user.Core
}

// Create adds a new tag to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var nt tag.NewTag
	if err := web.Decode(r, &nt); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if nt.Wid == "" {
		users, err := h.User.QueryByID(ctx, claims.Subject)
		if err != nil {
			return fmt.Errorf("unable to querying user: %w", err)
		}
		nt.Wid = users.DefaultWid
	} else {
		workspaces, err := h.Workspace.QueryByID(ctx, nt.Wid)
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
	}

	tags, err := h.Tag.Create(ctx, nt, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("tag[%+v]: %w", &tags, err)
		}
	}
	return web.Respond(ctx, w, tags, http.StatusCreated)
}

// Update updates a tag in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ug tag.UpdateTag
	if err := web.Decode(r, &ug); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	tagID := web.Param(r, "id")

	tags, err := h.Tag.QueryByID(ctx, tagID)
	if err != nil {
		switch {
		case errors.Is(err, tag.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, tag.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", tagID, err)
		}
	}
	workspaces, err := h.Workspace.QueryByID(ctx, tags.Wid)
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

	if err := h.Tag.Update(ctx, tagID, ug, v.Now); err != nil {
		switch {
		case errors.Is(err, tag.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, tag.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] tag[%+v]: %w", tagID, &ug, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a tag from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	tagID := web.Param(r, "id")

	tags, err := h.Tag.QueryByID(ctx, tagID)
	if err != nil {
		switch {
		case errors.Is(err, tag.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, tag.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", tagID, err)
		}
	}
	workspaces, err := h.Workspace.QueryByID(ctx, tags.Wid)
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

	if err := h.Tag.Delete(ctx, tagID); err != nil {
		switch {
		case errors.Is(err, tag.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", tagID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Package teamgrp maintains the group of handlers for projectuser access.
package teamgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/AhmedShaef/wakt/business/core/team"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
)

// Handlers manages the set of projectuser endpoints.
type Handlers struct {
	Team      team.Core
	Workspace workspace.Core
	User      user.Core
}

// Add adds a new team to the system.
func (h Handlers) Add(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var npu team.NewTeam
	if err := web.Decode(r, &npu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if npu.WID == "" {
		users, err := h.User.QueryByID(ctx, claims.Subject)
		if err != nil {
			return fmt.Errorf("unable to querying user: %w", err)
		}
		npu.WID = users.DefaultWid
	}

	projectUser, err := h.Team.QueryByID(ctx, npu.Puis)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, team.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying Team[%s]: %w", projectUser.ID, err)
		}
	}

	// If you are not an admin and looking to update a team you don't own.
	if !projectUser.Manager || npu.UID != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	clint, err := h.Team.Create(ctx, npu, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, team.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("team[%+v]: %w", &clint, err)
		}
	}

	return web.Respond(ctx, w, clint, http.StatusCreated)
}

// BulkUpdate updates a team in the system.
func (h Handlers) BulkUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd team.UpdateTeam
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	projctUserID := web.Param(r, "id")
	projectUserIDs := strings.Split(projctUserID, ",")

	for _, projectUserID := range projectUserIDs {
		projectUsers, err := h.Team.QueryByID(ctx, projectUserID)
		if err != nil {
			switch {
			case errors.Is(err, team.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, team.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", projectUserID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if !projectUsers.Manager || claims.Subject != projectUsers.UID {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.Team.Update(ctx, projectUserID, upd, v.Now); err != nil {
			switch {
			case errors.Is(err, team.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, team.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("ID[%s] Team[%+v]: %w", projectUserID, &upd, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// BulkDelete removes a team from the system.
func (h Handlers) BulkDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	projctUserID := web.Param(r, "id")
	projectUserIDs := strings.Split(projctUserID, ",")

	for _, projectUserID := range projectUserIDs {
		projectUsers, err := h.Team.QueryByID(ctx, projectUserID)
		if err != nil {
			switch {
			case errors.Is(err, team.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, team.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", projectUserID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if !projectUsers.Manager || claims.Subject != projectUsers.UID {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.Team.Delete(ctx, projectUserID); err != nil {
			switch {
			case errors.Is(err, team.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			default:
				return fmt.Errorf("ID[%s]: %w", projectUserID, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

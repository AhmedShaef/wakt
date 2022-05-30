// Package taskgrp maintains the group of handlers for task access.
package taskgrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/task"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
	"strings"
)

// Handlers manages the set of task endpoints.
type Handlers struct {
	Task      task.Core
	Workspace workspace.Core
}

// Create adds a new task to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var nt task.NewTask
	if err := web.Decode(r, &nt); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaces, err := h.Workspace.QueryByID(ctx, nt.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", nt.Wid, err)
		}
	}

	// If you are not an admin and looking to update a client you don't own.
	if workspaces.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	tsk, err := h.Task.Create(ctx, workspaces.Uid, nt, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, task.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("task[%+v]: %w", &tsk, err)
		}
	}

	return web.Respond(ctx, w, tsk, http.StatusCreated)
}

// QueryByID returns a task by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	taskID := web.Param(r, "id")

	tasks, err := h.Task.QueryByID(ctx, taskID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, task.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", taskID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != tasks.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	tsk, err := h.Task.QueryByID(ctx, taskID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, task.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", taskID, err)
		}
	}

	return web.Respond(ctx, w, tsk, http.StatusOK)
}

// BulkUpdate updates a task in the system.
func (h Handlers) BulkUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ut task.UpdateTask
	if err := web.Decode(r, &ut); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	tskID := web.Param(r, "id")
	taskIDs := strings.Split(tskID, ",")

	for _, taskID := range taskIDs {
		tasks, err := h.Task.QueryByID(ctx, taskID)
		if err != nil {
			switch {
			case errors.Is(err, task.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, task.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", taskID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if claims.Subject != tasks.Uid {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.Task.Update(ctx, taskID, ut, v.Now); err != nil {
			switch {
			case errors.Is(err, task.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, task.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("ID[%s] task[%+v]: %w", taskID, &ut, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// BulkDelete removes a task from the system.
func (h Handlers) BulkDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	tskID := web.Param(r, "id")
	taskIDs := strings.Split(tskID, ",")

	for _, taskID := range taskIDs {
		tasks, err := h.Task.QueryByID(ctx, taskID)
		if err != nil {
			switch {
			case errors.Is(err, task.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, task.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", taskID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if claims.Subject != tasks.Uid {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.Task.Delete(ctx, taskID); err != nil {
			switch {
			case errors.Is(err, task.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, task.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("ID[%s]: %w", taskID, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

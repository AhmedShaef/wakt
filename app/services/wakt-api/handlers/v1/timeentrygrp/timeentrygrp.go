// Package timeentrygrp maintains the group of handlers for timeEntry access.
package timeentrygrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/time_entry"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Handlers manages the set of timeEntry endpoints.
type Handlers struct {
	TimeEntry time_entry.Core
	Workspace workspace.Core
}

// Create adds a new timeEntry to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var nte time_entry.NewTimeEntry
	if err := web.Decode(r, &nte); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaces, err := h.Workspace.QueryByID(ctx, nte.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", nte.Wid, err)
		}
	}

	// If you are not an admin and looking to update a time entry you don't own.
	if workspaces.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	usr, err := h.TimeEntry.Create(ctx, nte, workspaces.Uid, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("timeEntry[%+v]: %w", &usr, err)
		}
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Start adds a new timeEntry to the system.
func (h Handlers) Start(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ste time_entry.StartTimeEntry
	if err := web.Decode(r, &ste); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	workspaces, err := h.Workspace.QueryByID(ctx, ste.Wid)
	if err != nil {
		switch {
		case errors.Is(err, workspace.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, workspace.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", ste.Wid, err)
		}
	}

	// If you are not an admin and looking to update a time entry you don't own.
	if workspaces.Uid != claims.Subject {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	usr, err := h.TimeEntry.Start(ctx, ste, workspaces.Uid, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("timeEntry[%+v]: %w", &usr, err)
		}
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Stop updates a timeEntry in the system.
func (h Handlers) Stop(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	timeEntryID := web.Param(r, "id")

	timeEntrys, err := h.TimeEntry.QueryByID(ctx, timeEntryID)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", timeEntryID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != timeEntrys.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	timeEntry, err := h.TimeEntry.Stop(ctx, timeEntryID, v.Now)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] time entry: %w", timeEntryID, err)
		}
	}

	return web.Respond(ctx, w, timeEntry, http.StatusNoContent)
}

// Update updates a timeEntry in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ute time_entry.UpdateTimeEntry
	if err := web.Decode(r, &ute); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	timeEntryID := web.Param(r, "id")

	timeEntrys, err := h.TimeEntry.QueryByID(ctx, timeEntryID)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", timeEntryID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != timeEntrys.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.TimeEntry.Update(ctx, timeEntryID, ute, v.Now); err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] time entry[%+v]: %w", timeEntryID, &ute, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a timeEntry from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	timeEntryID := web.Param(r, "id")

	timeEntrys, err := h.TimeEntry.QueryByID(ctx, timeEntryID)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", timeEntryID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != timeEntrys.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.TimeEntry.Delete(ctx, timeEntryID); err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", timeEntryID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// QueryByID returns a timeEntry by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	timeEntryID := web.Param(r, "id")

	timeEntrys, err := h.TimeEntry.QueryByID(ctx, timeEntryID)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying workspace[%s]: %w", timeEntryID, err)
		}
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if claims.Subject != timeEntrys.Uid {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	usr, err := h.TimeEntry.QueryByID(ctx, timeEntryID)
	if err != nil {
		switch {
		case errors.Is(err, time_entry.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, time_entry.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", timeEntryID, err)
		}
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

// QueryRunning returns a list of time entries with paging.
func (h Handlers) QueryRunning(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	timentry, err := h.TimeEntry.QueryRunning(ctx, claims.Subject, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for time entries: %w", err)
	}

	return web.Respond(ctx, w, timentry, http.StatusOK)
}

// QueryRange returns a list of time entries with paging.
func (h Handlers) QueryRange(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	start, err := time.Parse("2006-01-02T15:04:05+20:00", r.URL.Query().Get("start_date"))
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid start_date format, start_date[%s]", start), http.StatusBadRequest)
	}
	end, err := time.Parse("2006-01-02T15:04:05+20:00", r.URL.Query().Get("end_date"))
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid end_date format, end_date[%s]", end), http.StatusBadRequest)
	}

	timentry, err := h.TimeEntry.QueryRange(ctx, claims.Subject, pageNumber, rowsPerPage, start, end)
	if err != nil {
		return fmt.Errorf("unable to query for time entries: %w", err)
	}

	return web.Respond(ctx, w, timentry, http.StatusOK)
}

// UpdateTags updates a timeEntry in the system.
func (h Handlers) UpdateTags(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var ut time_entry.UpdateTimeEntryTags
	if err := web.Decode(r, &ut); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	timentryID := web.Param(r, "id")
	timeEntryIDs := strings.Split(timentryID, ",")

	for _, timeEntryID := range timeEntryIDs {
		timeEntries, err := h.TimeEntry.QueryByID(ctx, timeEntryID)
		if err != nil {
			switch {
			case errors.Is(err, time_entry.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, time_entry.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("querying workspace[%s]: %w", timeEntryID, err)
			}
		}

		// If you are not an admin and looking to retrieve someone other than yourself.
		if claims.Subject != timeEntries.Uid {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}

		if err := h.TimeEntry.UpdateTags(ctx, timeEntryID, ut, v.Now); err != nil {
			switch {
			case errors.Is(err, time_entry.ErrInvalidID):
				return v1Web.NewRequestError(err, http.StatusBadRequest)
			case errors.Is(err, time_entry.ErrNotFound):
				return v1Web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("ID[%s] time entry[%+v]: %w", timeEntryID, &ut, err)
			}
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// QueryDash returns a list of time entries with paging.
func (h Handlers) QueryDash(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	timentry, err := h.TimeEntry.QueryDash(ctx, claims.Subject)
	if err != nil {
		return fmt.Errorf("unable to query for time entries: %w", err)
	}

	return web.Respond(ctx, w, timentry, http.StatusOK)
}
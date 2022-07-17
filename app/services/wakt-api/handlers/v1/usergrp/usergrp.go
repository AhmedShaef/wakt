// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/business/core/project"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/core/workspace_user"
	feed "github.com/AhmedShaef/wakt/business/feed/geolocation"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/upload"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
	"strconv"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	User          user.Core
	Workspace     workspace.Core
	WorkspaceUser workspace_user.Core
	Auth          *auth.Auth
	Project       project.Core
}

// SignUp adds a new user to the system.
func (h Handlers) SignUp(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	inviteKey := web.Param(r, "invite_key")
	var inviteUID string
	var inviteWID string

	if inviteKey != "" {
		inviteUID = inviteKey[9:54]
		inviteWID = inviteKey[55:91]
	}

	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	usr, err := h.User.Create(ctx, nu, v.Now, inviteUID, inviteWID)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return v1Web.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}
	nw := workspace.NewWorkspace{
		Name: usr.FullName,
		UID:  usr.ID,
	}
	initWorkspace, err := h.Workspace.Create(ctx, nw, v.Now)
	if err != nil {
		return fmt.Errorf("initWorkspace[%+v]: %w", &initWorkspace, err)
	}

	geo, err := feed.GetGeolocation(r.RemoteAddr)
	if err != nil {
		return fmt.Errorf("unable to get geolocation: %w", err)
	}

	uu := user.UpdateUser{
		DefaultWid: &initWorkspace.ID,
		TimeZone:   &geo.TimeZone,
	}

	err = h.User.Update(ctx, usr.ID, uu, v.Now)
	if err != nil {
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}

	finalUser, err := h.User.QueryByID(ctx, usr.ID)
	if err != nil {
		return fmt.Errorf("finalUser[%+v]: %w", &finalUser, err)
	}

	return web.Respond(ctx, w, finalUser, http.StatusCreated)
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

	var upd user.UpdateUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if err := h.User.Update(ctx, claims.Subject, upd, v.Now); err != nil {
		return fmt.Errorf("ID[%s] User[%+v]: %w", claims.Subject, &upd, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// UpdateImage updates a user in the system.
func (h Handlers) UpdateImage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	file, handler, err := r.FormFile("profileImage")
	if err != nil {
		return fmt.Errorf("unable to get file %v: %w", handler.Filename, err)
	}
	defer file.Close()

	name, err := upload.Image(file)
	if err != nil {
		return fmt.Errorf("unable to upload image: %w", err)
	}
	upd := user.UpdateImage{
		ImageName: name,
	}

	if err := h.User.UpdateImage(ctx, claims.Subject, upd, v.Now); err != nil {
		return fmt.Errorf("ID[%s] User[%+v]: %w", claims.Subject, &upd, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// ChangePassword updates a user in the system.
func (h Handlers) ChangePassword(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var chp user.ChangePassword
	if err := web.Decode(r, &chp); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if err := h.User.ChangePassword(ctx, claims.Subject, chp, v.Now); err != nil {
		return fmt.Errorf("ID[%s] User[%+v]: %w", claims.Subject, &chp, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// QueryByID returns a user by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	usr, err := h.User.QueryByID(ctx, claims.Subject)
	if err != nil {
		return fmt.Errorf("ID[%s]: %w", claims.Subject, err)
	}
	// TODO: add related date
	return web.Respond(ctx, w, usr, http.StatusOK)
}

// NewToken provides an API token for the authenticated user.
func (h Handlers) NewToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return v1Web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := h.User.Authenticate(ctx, email, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			return v1Web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authenticating: %w", err)
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = h.Auth.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}

// QueryUserProjects returns a list of projects for the user.
func (h Handlers) QueryUserProjects(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

	UserProject, err := h.Project.QueryUserProjects(ctx, claims.Subject, pageNumber, rowsPerPage)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, user.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("querying user project[%s]: %w", claims.Subject, err)
		}
	}

	return web.Respond(ctx, w, UserProject, http.StatusOK)
}

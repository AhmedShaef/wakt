// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/AhmedShaef/wakt/app/tooling/geolocation"
	"github.com/AhmedShaef/wakt/app/tooling/uploader"
	"github.com/AhmedShaef/wakt/business/core/user"
	"github.com/AhmedShaef/wakt/business/core/workspace"
	"github.com/AhmedShaef/wakt/business/core/workspace_user"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	v1Web "github.com/AhmedShaef/wakt/business/web/v1"
	"github.com/AhmedShaef/wakt/foundation/web"
	"net/http"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	User          user.Core
	Workspace     workspace.Core
	WorkspaceUser workspace_user.Core
	Auth          *auth.Auth
	APIKey        *geolocation.Config
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
		Uid:  usr.ID,
	}
	initWorkspace, err := h.Workspace.Create(ctx, nw, v.Now)
	if err != nil {
		return fmt.Errorf("initWorkspace[%+v]: %w", &initWorkspace, err)
	}

	geo, err := geolocation.GetGeolocation(r.RemoteAddr, *h.APIKey)
	if err != nil {

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

	me := web.Param(r, "id")
	if me != "me" {
		return v1Web.NewRequestError(err, http.StatusBadRequest)
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
		return fmt.Errorf("unable to get file %S: %w", handler.Filename, err)
	}
	defer file.Close()

	name, err := uploader.UploadImage(file)
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

	me := web.Param(r, "id")
	if me != "me" {
		return v1Web.NewRequestError(err, http.StatusBadRequest)
	}
	usr, err := h.User.QueryByID(ctx, claims.Subject)
	if err != nil {
		return fmt.Errorf("ID[%s]: %w", claims.Subject, err)
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

// NewToken provides an API token for the authenticated user.
func (h Handlers) NewToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return v1Web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := h.User.Authenticate(ctx, v.Now, email, pass)
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

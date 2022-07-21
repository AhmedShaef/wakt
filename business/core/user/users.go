// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AhmedShaef/wakt/business/core/user/db"
	"github.com/AhmedShaef/wakt/business/sys/auth"
	"github.com/AhmedShaef/wakt/business/sys/database"
	"github.com/AhmedShaef/wakt/business/sys/validate"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrInvalidEmail          = errors.New("email is not valid")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrInvalidPassword       = errors.New("password is not valid")
)

// Core manages the set of APIs for user access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu NewUser, now time.Time, inviteDate ...string) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	var iuID string
	var iwID string
	if len(inviteDate) == 2 {
		if err := validate.CheckID(inviteDate[0]); err == nil {
			iuID = inviteDate[0]
		}
		if err := validate.CheckID(inviteDate[1]); err != nil {
			iwID = inviteDate[1]
		}
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	name := strings.Split(nu.Email, "@")

	var userID string
	var workspaceID string
	if iuID != "" {
		userID = iuID
	} else {
		userID = validate.GenerateID()
	}

	if iwID != "" {
		workspaceID = iwID
	} else {
		workspaceID = ""
	}

	dbusr := db.User{
		ID:              userID,
		DefaultWid:      validate.GenerateID(),
		Email:           nu.Email,
		PasswordHash:    hashPassword,
		FullName:        name[0],
		TimeOfDayFormat: "h:mm tt",
		DateFormat:      "MM/DD/YYYY",
		BeginningOfWeek: 1,
		Language:        "en_US",
		ImageURL:        "",
		DateCreated:     now,
		DateUpdated:     now,
		TimeZone:        "",
		Invitation:      []string{workspaceID},
		DurationFormat:  "",
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		if err := c.store.Tran(tx).Create(ctx, dbusr); err != nil {
			if errors.Is(err, database.ErrDBDuplicatedEntry) {
				return fmt.Errorf("create: %w", ErrUniqueEmail)
			}
			return fmt.Errorf("create: %w", err)
		}
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return User{}, fmt.Errorf("tran: %w", err)
	}

	return toUser(dbusr), nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, userID string, uu UpdateUser, now time.Time) error {
	if err := validate.CheckID(userID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbuser, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}

	if uu.DefaultWid != nil {
		dbuser.DefaultWid = *uu.DefaultWid
	}

	if uu.Email != nil {
		dbuser.Email = *uu.Email
	}
	if uu.FullName != nil {
		dbuser.FullName = *uu.FullName
	}
	if uu.TimeOfDayFormat != nil {
		dbuser.TimeOfDayFormat = *uu.TimeOfDayFormat
	}
	if uu.DateFormat != nil {
		dbuser.DateFormat = *uu.DateFormat
	}
	if uu.BeginningOfWeek != nil {
		dbuser.BeginningOfWeek = *uu.BeginningOfWeek
	}
	if uu.Language != nil {
		dbuser.Language = *uu.Language
	}
	if uu.TimeZone != nil {
		dbuser.TimeZone = *uu.TimeZone
	}
	if uu.Invitation != nil {
		dbuser.Invitation = uu.Invitation
	}
	if uu.DurationFormat != nil {
		dbuser.DurationFormat = *uu.DurationFormat
	}
	dbuser.DateUpdated = now

	if err := c.store.Update(ctx, dbuser); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// ChangePassword replaces a user document in the database.
func (c Core) ChangePassword(ctx context.Context, userID string, cp ChangePassword, now time.Time) error {
	if err := validate.CheckID(userID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(cp); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbuser, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}

	if err := bcrypt.CompareHashAndPassword(dbuser.PasswordHash, []byte(cp.OldPassword)); err != nil {
		return ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(*cp.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generating password hash: %w", err)
	}
	dbuser.PasswordHash = hash
	dbuser.DateUpdated = now

	if err := c.store.Update(ctx, dbuser); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// UpdateImage replaces a user document in the database.
func (c Core) UpdateImage(ctx context.Context, userID string, ui UpdateImage, now time.Time) error {
	if err := validate.CheckID(userID); err != nil {
		return ErrInvalidID
	}

	dbuser, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}
	if ui.ImageName != "" {
		dbuser.ImageURL = "app/foundation/upload/assets" + ui.ImageName
	}
	dbuser.DateUpdated = now

	if err := c.store.Update(ctx, dbuser); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// QueryByID gets the specified user from the database.
func (c Core) QueryByID(ctx context.Context, userID string) (User, error) {
	if err := validate.CheckID(userID); err != nil {
		return User{}, ErrInvalidID
	}

	dbuser, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("query: %w", err)
	}

	return toUser(dbuser), nil
}

// QueryByEmail gets the specified user from the database by email.
func (c Core) QueryByEmail(ctx context.Context, email string) (User, error) {

	// Email Validate function in validate.
	if !validate.CheckEmail(email) {
		return User{}, ErrInvalidEmail
	}

	dbUsr, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("query: %w", err)
	}

	return toUser(dbUsr), nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, email, password string) (auth.Claims, error) {

	// Email Validate function in validate.
	if !validate.CheckEmail(email) {
		return auth.Claims{}, ErrInvalidEmail
	}

	dbUser, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return auth.Claims{}, ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("query: %w", err)
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(dbUser.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUser.ID,
			Issuer:    "wakt project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
	return claims, nil
}

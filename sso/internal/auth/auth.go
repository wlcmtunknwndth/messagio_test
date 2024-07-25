package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/jwt"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/domain/models"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const scope = "sso.inner.auth."

type Auth struct {
	log      *slog.Logger
	storage  Storage
	tokenTTL time.Duration
}

func (a *Auth) GetTokenTTL() time.Duration {
	return a.tokenTTL
}

type Storage interface {
	User(ctx context.Context, email string) (*models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	SaveUser(ctx context.Context, user string, passHash []byte) (int64, error)
}

func New(log *slog.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, username, password string) (string, error) {
	const op = scope + "Login"

	usr, err := a.storage.User(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, internal.ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(usr.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// admin check
	isAdmin, err := a.IsAdmin(ctx, usr.ID)
	if err != nil {
		isAdmin = false
	}

	token, err := jwt.NewToken(usr.ID, usr.Username, a.tokenTTL, isAdmin)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, username, pass string) (int64, error) {
	const op = scope + ".RegisterNewUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.storage.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {

			return 0, fmt.Errorf("%s: %w", op, internal.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = scope + "IsAdmin"

	isAdmin, err := a.storage.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return false, fmt.Errorf("%s: %w", op, internal.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

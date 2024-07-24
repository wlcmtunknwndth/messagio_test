package authHandler

import (
	"context"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"log/slog"
	"net/http"
)

type AuthService interface {
	User(ctx context.Context, email string) (*api.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	SaveUser(ctx context.Context, user string, passHash []byte) (int64, error)
}

type Auth struct {
	service AuthService
	log     *slog.Logger
}

const (
	scope = "sso.internal.handlers.authHandler."

	errUserNotFound        = "User not found"
	errInvalidCredentials  = "Invalid credentials"
	errInternalServerError = "Internal server error"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	const op = scope + "Login"

	var usr api.User

}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	const op = scope + "Register"

}

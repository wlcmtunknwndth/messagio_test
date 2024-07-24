package authHandler

import (
	"context"
	"encoding/json"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"github.com/wlcmtunknwndth/messagio_test/common/httpResp"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	RegisterNewUser(ctx context.Context, email, pass string) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	GetTokenTTL() time.Duration
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
	errBadRequest          = "Bad request"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	const op = scope + "Login"

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			a.log.Error("couldn't close request body", sl.Op(op), sl.Err(err))
			return
		}
		return
	}(r.Body)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("couldn't read body", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, errBadRequest)
		return
	}

	var usr api.User
	if err = json.Unmarshal(data, &usr); err != nil {
		a.log.Error("couldn't unmarshal request", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := a.service.Login(r.Context(), usr.Username, usr.Password)
	if err != nil {
		a.log.Error("couldn't login user", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	httpResp.WriteToken(w, token, a.service.GetTokenTTL())
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	const op = scope + "Register"

}

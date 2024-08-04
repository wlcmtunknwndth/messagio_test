package authHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/httpResp"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/domain/models"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	RegisterNewUser(ctx context.Context, username, pass string) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	GetTokenTTL() time.Duration
}

type Auth struct {
	service AuthService
	log     *slog.Logger
}

func New(service AuthService, log *slog.Logger) *Auth {
	return &Auth{
		service: service,
		log:     log,
	}
}

const (
	scope = "sso.internal.handlers.authHandler."

	errUserNotFound        = "User not found"
	errUserExists          = "User exits"
	errInvalidCredentials  = "Invalid credentials"
	errInternalServerError = "Internal server error"
	errBadRequest          = "Bad request"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	const op = scope + "Login"

	usr, err := extractUserCredentials(r)
	if err != nil {
		a.log.Error("couldn't handle request body", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, errBadRequest)
		return
	}

	token, err := a.service.Login(r.Context(), usr.Username, usr.Password)
	if err != nil {
		if errors.Is(err, internal.ErrUserNotFound) {
			a.log.Error("user exists", sl.Op(op), sl.Err(err))
			httpResp.Write(w, http.StatusNotFound, errUserNotFound)
			return
		}
		if errors.Is(err, internal.ErrInvalidCredentials) {
			a.log.Error("invalid credentials", sl.Op(op), sl.Err(err))
			httpResp.Write(w, http.StatusForbidden, errInvalidCredentials)
			return
		}
		a.log.Error("couldn't login user", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	httpResp.WriteToken(w, token, a.service.GetTokenTTL())
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	const op = scope + "Register"

	usr, err := extractUserCredentials(r)
	if err != nil {
		a.log.Error("couldn't handle request body", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, errBadRequest)
		return
	}

	id, err := a.service.RegisterNewUser(r.Context(), usr.Username, usr.Password)
	if err != nil {
		if errors.Is(err, internal.ErrUserExists) {
			a.log.Error("user already exists", sl.Op(op), sl.Err(err))
			httpResp.Write(w, http.StatusForbidden, errUserExists)
			return
		}
		a.log.Error("couldn't register new user", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, errInternalServerError)
		return
	}

	httpResp.Write(w, http.StatusCreated, fmt.Sprintf("Created user: id=%d", id))
}

//func (a *Auth) IsAdmin(w http.ResponseWriter, r *http.Request) {
//	const op = scope + "IsAdmin"
//
//	idQry, ok := r.URL.Query()["id"]
//	if !ok {
//		a.log.Error("couldn't find id in query", sl.Op(op))
//		httpResp.Write(w, http.StatusNotFound, errIdQuery)
//		return
//	}
//
//	id, err := strconv.ParseInt(idQry[0], 10, 64)
//	if err != nil {
//		a.log.Error("couldn't cast string to int64", sl.Op(op), sl.Err(err))
//		httpResp.Write(w, http.StatusBadRequest, errBadRequest)
//		return
//	}
//
//	res, err := a.handler.IsAdmin(r.Context(), id)
//	if err != nil {
//		a.log.Error("couldn't determine if user is admin", sl.Op(op), sl.Err(err))
//		httpResp.Write(w, http.StatusInternalServerError, errInternalServerError)
//		return
//	}
//
//	if res {
//		httpResp.Write(w, http.StatusOK, "ok")
//	}else{
//		httpResp.Write(w, http.)
//	}
//
//}

func extractUserCredentials(r *http.Request) (*models.UserAPI, error) {
	const op = scope + "extractUserCredentials"

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var usr models.UserAPI
	if err = json.Unmarshal(data, &usr); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &usr, nil
}

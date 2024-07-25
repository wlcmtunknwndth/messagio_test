package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

const (
	scope    = "sso.internal.app."
	register = "register"
	login    = "login"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type App struct {
	handler AuthHandler
	log     *slog.Logger
	server  *http.Server
}

func New(Address string, timeout, idleTimeout time.Duration, auth AuthHandler, log *slog.Logger) *App {
	router := chi.NewRouter()
	router.Post(login, auth.Login)
	router.Post(register, auth.Register)
	return &App{
		handler: auth,
		log:     nil,
		server: &http.Server{
			Addr:         Address,
			WriteTimeout: timeout,
			IdleTimeout:  idleTimeout,
			Handler:      router,
		},
	}
}

func (a *App) Run() error {
	const op = scope + "Run"

	if err := a.server.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Close() error {
	const op = scope + "Close"

	if err := a.server.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

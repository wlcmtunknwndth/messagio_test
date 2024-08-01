package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

const (
	scope    = "sso.internal.app."
	register = "/register"
	login    = "/login"
)

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type App struct {
	handler AuthHandler
	server  *http.Server
}

func New(Address string, timeout, idleTimeout time.Duration, auth AuthHandler) *App {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post(login, auth.Login)
	router.Post(register, auth.Register)
	return &App{
		handler: auth,
		server: &http.Server{
			Addr:         Address,
			WriteTimeout: timeout,
			ReadTimeout:  timeout,
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

func (a *App) Close(ctx context.Context) error {
	const op = scope + "Close"

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

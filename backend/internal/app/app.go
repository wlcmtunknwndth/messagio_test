package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

const (
	scope = "backend.internal.app."

	getChats = "/chats"
	getChat  = "/chat"
	sendMsg  = "/send"
)

type MessagerHandler interface {
	HandleMessage(w http.ResponseWriter, r *http.Request)
	HandleChatRequest(w http.ResponseWriter, r *http.Request)
	GetChats(w http.ResponseWriter, r *http.Request)
}

type App struct {
	handler MessagerHandler
	server  *http.Server
}

func New(address string, timeout, idleTimeout time.Duration, handler MessagerHandler) *App {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post(sendMsg, handler.HandleMessage)
	router.Get(getChat, handler.HandleChatRequest)
	router.Get(getChats, handler.GetChats)

	return &App{
		handler: handler,
		server: &http.Server{
			Addr:         address,
			WriteTimeout: timeout,
			ReadTimeout:  timeout,
			IdleTimeout:  idleTimeout,
			Handler:      router,
		},
	}
}

func (a *App) Run() error {
	const op = scope + "Run"

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

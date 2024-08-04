package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type HandlerHTTP interface {
	MessagesReceivedByUser(w http.ResponseWriter, r *http.Request)
	MessagesReceived(w http.ResponseWriter, r *http.Request)
	MessagesSentByUser(w http.ResponseWriter, r *http.Request)
	Close() error
}

const scope = "stats.internal.app."

type Broker interface {
	FetchMessages(ctx context.Context)
	Close() error
}

type App struct {
	handler HandlerHTTP
	broker  Broker
	server  *http.Server
}

func New(handler HandlerHTTP, broker Broker, address string, writeTimeout time.Duration, idleTimeout time.Duration) *App {
	router := chi.NewRouter()

	router.Get("MessagesReceived", handler.MessagesReceived)
	router.Get("MessagesReceivedByUser", handler.MessagesReceivedByUser)
	router.Get("MessagesSentByUser", handler.MessagesSentByUser)
	return &App{
		handler: handler,
		broker:  broker,
		server: &http.Server{
			Addr:         address,
			WriteTimeout: writeTimeout,
			ReadTimeout:  writeTimeout,
			IdleTimeout:  idleTimeout,
			Handler:      router,
		},
	}
}

func (a *App) Start(ctx context.Context) error {
	const op = scope + "Start"

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go a.broker.FetchMessages(ctx)

	if err := a.server.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) GracefulShutdown(ctx context.Context) error {
	const op = scope + "GracefulShutdown"

	err := errors.Join(a.handler.Close(), a.broker.Close(), a.server.Shutdown(ctx))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

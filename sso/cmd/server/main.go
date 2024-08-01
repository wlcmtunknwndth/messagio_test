package main

import (
	"context"
	"github.com/wlcmtunknwndth/messagio_test/common/logger"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/app"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/auth"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/config"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/handlers/authHandler"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/storage/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const scope = "sso.cmd.server.main"

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Configuration loaded", slog.Any("config", cfg))

	storage, err := postgres.New(&cfg.DB)
	if err != nil {
		slog.Error("couldn't connect to postgres", sl.Op(scope), sl.Err(err))
		return
	}

	authService := auth.New(log, storage, cfg.TokenTTL)

	application := app.New(cfg.Server.Addr, cfg.Server.Timeout, cfg.Server.IdleTimeout, authHandler.New(authService, log))

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := application.Run(); err != nil {
			log.Error("couldn't run server", sl.Op(scope), sl.Err(err))
		}
		return
	}()

	slog.Info("server started", slog.String("Address", cfg.Server.Addr))

	<-stop
	if err = application.Close(context.Background()); err != nil {
		log.Error("couldn't close application", sl.Op(scope), sl.Err(err))
		return
	}
	return
}

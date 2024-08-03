package main

import (
	"context"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/app"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/broker/kafka"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/config"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/handlers/messageHandler"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/messager"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/storage/postgres"
	"github.com/wlcmtunknwndth/messagio_test/common/logger"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const scope = "backend.cmd.server.main"

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Configuration loaded", slog.Any("Config", cfg))

	storage, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Error("Failed to connect to database", sl.Op(scope), sl.Err(err))
		return
	}
	kfk, err := kafka.New(&cfg.Broker, log)
	if err != nil {
		log.Error("Failed to connect to kafka", sl.Op(scope), sl.Err(err))
		return
	}

	application := app.New(
		cfg.Server.Address, cfg.Server.Timeout, cfg.Server.IdleTimeout,
		messageHandler.New(messager.New(storage, kfk), log),
	)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Run(); err != nil {
			log.Error("Failed to start application", sl.Op(scope), sl.Err(err))
		}
		return
	}()

	slog.Info("Server started", slog.String("Address", cfg.Server.Address))

	<-stop

	if err = application.Close(context.Background()); err != nil {
		log.Error("couldn't close application", sl.Op(scope), sl.Err(err))
		return
	}
	return
}

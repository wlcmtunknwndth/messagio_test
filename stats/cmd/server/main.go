package main

import (
	"context"
	"github.com/wlcmtunknwndth/messagio_test/common/logger"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"github.com/wlcmtunknwndth/stats/internal/app"
	"github.com/wlcmtunknwndth/stats/internal/broker/kafka"
	"github.com/wlcmtunknwndth/stats/internal/config"
	"github.com/wlcmtunknwndth/stats/internal/handlers/RestAPI"
	"github.com/wlcmtunknwndth/stats/internal/storage/postgres"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const scope = "main"

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Config found", slog.Any("Config", cfg))

	storage, err := postgres.New(&cfg.DB)
	if err != nil {
		log.Error("Couldn't connect to postgres", sl.Op(scope), sl.Err(err))
		return
	}

	broker, err := kafka.New(&cfg.Broker, storage, log)
	if err != nil {
		log.Error("Couldn't connect to kafka", sl.Op(scope), sl.Err(err))
		return
	}

	application := app.New(RestAPI.New(storage, log), broker, cfg.Server.Address, cfg.Server.Timeout, cfg.Server.IdleTimeout)

	go func() {
		err := application.Start(context.Background())
		if err != nil {
			log.Error("Server closed", sl.Op(scope), sl.Err(err))
			return
		}
		return
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGTERM)

	<-ch
	if err = application.GracefulShutdown(context.Background()); err != nil {
		log.Error("Couldn't close server", sl.Op(scope), sl.Err(err))
		return
	}
	return
}

package server

import (
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

const scope = "main"

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	storage, err := postgres.New(&cfg.DB)
	if err != nil {
		slog.Error("couldn't connect to postgres", sl.Op(scope), sl.Err(err))
		return
	}

	authService := auth.New(log, storage, cfg.TokenTTL)

	application := app.New(cfg.Server.Addr, cfg.Server.Timeout, cfg.Server.IdleTimeout, authHandler.New(authService, log), log)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := application.Run(); err != nil {
			log.Error("couldn't run server", sl.Op(scope), sl.Err(err))
		}
		return
	}()

	<-stop
	if err = application.Close(); err != nil {
		slog.Error("couldn't close application", sl.Op(scope), sl.Err(err))
		return
	}
	return
}

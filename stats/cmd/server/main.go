package main

import (
	"github.com/wlcmtunknwndth/messagio_test/common/logger"
	"github.com/wlcmtunknwndth/stats/internal/config"
	"log/slog"
)

const scope = "main"

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("Config found", slog.Any("Config", cfg))
}

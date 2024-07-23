package main

import (
	"github.com/wlcmtunknwndth/messagio_test/common/config"
	"github.com/wlcmtunknwndth/messagio_test/common/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

}

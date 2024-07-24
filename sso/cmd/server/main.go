package server

import (
	"github.com/wlcmtunknwndth/messagio_test/common/logger"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

}

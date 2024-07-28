package postgres

import (
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/backend/internal/config"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

const (
	scope = "backend.internal.storage.postgres."
)

func New(cfg *config.DataBase) (*Storage, error) {
	const op = scope + "New"

	connStr := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=%s",
		cfg.DbUser, cfg.DbPass, cfg.Port,
		cfg.DbName, cfg.SslMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.AutoMigrate(&api.Message{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

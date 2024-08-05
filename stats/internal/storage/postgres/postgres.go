package postgres

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"github.com/wlcmtunknwndth/stats/internal/config"
	"github.com/wlcmtunknwndth/stats/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const scope = "stats.internal.storage.postgres."

type Storage struct {
	db *gorm.DB
}

func New(cfg *config.DataBase) (*Storage, error) {
	const op = scope + "New"

	connStr := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=%s",
		cfg.DbUser, cfg.DbPass, cfg.Port, cfg.DbName, cfg.SslMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.AutoMigrate(&models.MsgCount{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const op = scope + "Close"

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err = sqlDB.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CountMessage(ctx context.Context, message *api.Message) error {
	const op = scope + "CountMessage"

	if err := s.db.WithContext(ctx).Model(&models.MsgCount{}).Save(&models.MsgCount{
		UserID:    message.UserID,
		PalID:     message.PalID,
		CreatedAt: message.CreatedAt,
	}).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

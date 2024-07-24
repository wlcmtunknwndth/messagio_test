package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/config"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	scope = "sso.internal.storage.postgres."
)

type Storage struct {
	db *gorm.DB
}

func New(config *config.DataBase) (*Storage, error) {
	const op = scope + "New"

	connStr := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=%s",
		config.DbUser, config.DbPass, config.Port,
		config.DbName, config.SslMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.Admin{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

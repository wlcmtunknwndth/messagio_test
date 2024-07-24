package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/sso/internal/domain/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = scope + "User"

	var usr models.User
	if err := s.db.WithContext(ctx).First(&models.User{}).Where("email = ?", email).Scan(&usr); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err.Error)
	}

	return &usr, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = scope + "IsAdmin"

	var usr models.Admin
	if err := s.db.WithContext(ctx).First(&usr, userID).Scan(&usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return usr.IsAdmin, nil
}

func (s *Storage) SaveUser(ctx context.Context, user string, passHash []byte) (int64, error) {
	const op = scope + "SaveUser"

	var usr models.User
	usr.Username = user
	usr.PassHash = passHash
	if res := s.db.WithContext(ctx).Model(&models.User{}).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Create(&usr); res.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, res.Error)
	}

	return usr.ID, nil
}

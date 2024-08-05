package postgres

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/stats/internal/domain/models"
)

func (s *Storage) MessagesReceived(ctx context.Context, since, to int64) (*models.Stats, error) {
	const op = scope + "MessagesReceived"

	var counter int64
	if err := s.db.WithContext(ctx).Model(&models.MsgCount{}).
		Where("created_at > ? and created_at < ??", since, to).
		Count(&counter).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Stats{
		MessagesCounter: counter,
		Since:           since,
		To:              to,
	}, nil
}

func (s *Storage) MessagesSentByUser(ctx context.Context, userID, since, to int64) (*models.UserStats, error) {
	const op = scope + "MessagesSentByUser"

	var counter int64

	if err := s.db.WithContext(ctx).Model(&models.MsgCount{}).
		Where("user_id = ? and created_at > ? and created_at < ?", userID, since, to).
		Count(&counter).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.UserStats{
		UserID: userID,
		Stats: models.Stats{
			MessagesCounter: counter,
			Since:           since,
			To:              to,
		},
	}, nil
}

func (s *Storage) MessagesReceivedByUser(ctx context.Context, userID, since, to int64) (*models.UserStats, error) {
	const op = scope + "MessagesSentByUser"

	var counter int64
	if err := s.db.WithContext(ctx).Model(&models.MsgCount{}).
		Where("pal_id = ? and created_at > ? and created_at < ?", userID, since, to).
		Count(&counter).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.UserStats{
		UserID: userID,
		Stats: models.Stats{
			MessagesCounter: counter,
			Since:           since,
			To:              to,
		},
	}, nil
}

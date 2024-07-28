package postgres

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"gorm.io/gorm/clause"
)

func (s *Storage) Save(ctx context.Context, input *api.Message) (int64, error) {
	const op = scope + "Save"

	if err := s.db.WithContext(ctx).Model(&api.Message{}).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Save(input).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return input.ID, nil
}

func (s *Storage) ChatMessages(ctx context.Context, userID int64, palID int64, limit, offset int) ([]api.Message, error) {
	const op = scope + "ChatMessages"

	var chat []api.Message

	if err := s.db.WithContext(ctx).Where("user_id = ? AND pal_id = ? OR user_id = ? AND pal_id = ?", userID, palID, palID, userID).Order("created_at desc").Limit(limit).Offset(offset).Find(&chat).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (s *Storage) Chats(ctx context.Context, userID int64) ([]api.Message, error) {
	const op = scope + "Chats"

	var messages []api.Message

	subQuery := s.db.WithContext(ctx).Model(&api.Message{}).
		Select("pal_id, MAX(created_at) AS last_created_at").
		Group("pal_id")

	if err := s.db.WithContext(ctx).Model(&api.Message{}).
		Joins("JOIN (?) AS last_messages ON messages.pal_id = last_messages.pal_id AND messages.created_at = last_messages.last_created_at", subQuery).
		Scan(&messages).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}

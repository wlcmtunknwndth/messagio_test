package messager

import (
	"context"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
)

type Storage interface {
	Save(ctx context.Context, input *api.Message) (int64, error)
	ChatMessages(ctx context.Context, id int64, palID int64, limit, offset int) ([]api.Message, error)
	Chats(ctx context.Context, id int64) ([]api.Message, error)
}

type Broker interface {
	CountMessageSent(ctx context.Context, id, palID, createdAtUnix int64) error
}

const scope = "backend.internal.messager."

type Messager struct {
	storage Storage
	broker  Broker
}

func (m *Messager) SaveMessage(ctx context.Context, input *api.Message) (int64, error) {
	const op = scope + "SaveMessage"

	id, err := m.storage.Save(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err = m.broker.CountMessageSent(ctx, input.UserID, input.PalID, input.CreatedAt); err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (m *Messager) GetChat(ctx context.Context, id, palID int64, limit, offset int) ([]api.Message, error) {
	const op = scope + "GetChat"

	chat, err := m.storage.ChatMessages(ctx, id, palID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chat, nil
}

func (m *Messager) GetChats(ctx context.Context, id int64) ([]api.Message, error) {
	const op = scope + "GetChats"

	chats, err := m.storage.Chats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chats, nil
}

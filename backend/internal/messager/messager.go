package messager

import (
	"context"
	"errors"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
)

type Storage interface {
	Save(ctx context.Context, input *api.Message) (int64, error)
	ChatMessages(ctx context.Context, id int64, palID int64, limit, offset int) ([]api.Message, error)
	Chats(ctx context.Context, id int64) ([]api.Message, error)
	Close() error
}

type Broker interface {
	CountMessageSent(ctx context.Context, msg *api.Message) error
	Close() error
}

const scope = "backend.internal.messager."

type Messager struct {
	storage Storage
	broker  Broker
}

func New(storage Storage, broker Broker) *Messager {
	return &Messager{
		storage: storage,
		broker:  broker,
	}
}

func (m *Messager) SaveMessage(ctx context.Context, input *api.Message) (int64, error) {
	const op = scope + "SaveMessage"

	id, err := m.storage.Save(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	input.ID = id
	if err = m.broker.CountMessageSent(ctx, input); err != nil {
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

func (m *Messager) Close() error {
	const op = scope + "Close"

	err := errors.Join(m.broker.Close(), m.storage.Close())

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

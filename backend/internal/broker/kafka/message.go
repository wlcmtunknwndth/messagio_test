package kafka

import (
	"context"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
)

type Storage interface {
	Chats(ctx context.Context, id int64) []api.Message
	Send(ctx context.Context, id int64, palID int64) error
}

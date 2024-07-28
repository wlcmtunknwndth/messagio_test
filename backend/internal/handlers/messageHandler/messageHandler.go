package messageHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/domain/api"
	"github.com/wlcmtunknwndth/messagio_test/common/httpResp"
	"github.com/wlcmtunknwndth/messagio_test/common/jwt"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"io"
	"log/slog"
	"net/http"
)

type Messager interface {
	SaveMessage(ctx context.Context, input *api.Message) (int64, error)
	GetChat(ctx context.Context, id, palID int64, limit, offset int) ([]api.Message, error)
	GetChats(ctx context.Context, id int64) ([]api.Message, error)
}

type Handler struct {
	log      *slog.Logger
	messager Messager
}

const (
	badRequest          = "Bad request"
	internalServerError = "Internal server error"
	palID               = "pal_id"
	offset              = "offset"
	limit               = "limit"
)

const scope = "backend.internal.messageHandler."

func (h *Handler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	const op = scope + "HandleMessage"

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		httpResp.Write(w, http.StatusBadRequest, badRequest)
		h.log.Error("couldn't read body", sl.Op(op), sl.Err(err))
		return
	}

	var input api.Input
	if err = json.Unmarshal(data, &input); err != nil {
		httpResp.Write(w, http.StatusBadRequest, badRequest)
		h.log.Error("couldn't unmarshal data", sl.Op(op), sl.Err(err))
		return
	}

	jwt.GetID()

	var msg api.Message

	id, err := h.messager.SaveMessage(r.Context())
	if err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't save message", sl.Op(op), sl.Err(err))
		return
	}

	httpResp.Write(w, http.StatusCreated, fmt.Sprintf("created message: id = %d", id))
}

func (h *Handler) HandleChatRequest(w http.ResponseWriter, r *http.Request) {
	const op = scope + "HandleChatRequest"

	r.URL.Query()[]

}

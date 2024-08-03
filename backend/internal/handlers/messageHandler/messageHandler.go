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
	"strconv"
	"time"
)

type Messager interface {
	SaveMessage(ctx context.Context, input *api.Message) (int64, error)
	GetChat(ctx context.Context, id, palID int64, limit, offset int) ([]api.Message, error)
	GetChats(ctx context.Context, id int64) ([]api.Message, error)
	Close() error
}

type HandlerHTTP struct {
	log      *slog.Logger
	messager Messager
}

func New(msgr Messager, log *slog.Logger) *HandlerHTTP {
	return &HandlerHTTP{
		log:      log,
		messager: msgr,
	}
}

const (
	badRequest          = "Bad request"
	internalServerError = "Internal server error"
	tokenNotFound       = "Token not found"
	unauthorized        = "Unauthorized"
	queryNotFound       = "Query key/values not found"
	queryBadValue       = "Query bad value"

	palIdKey  = "pal_id"
	offsetKey = "offset"
	limitKey  = "limit"
)

const scope = "backend.internal.messageHandler."

func (h *HandlerHTTP) HandleMessage(w http.ResponseWriter, r *http.Request) {
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

	token, err := httpResp.GetToken(r)
	if err != nil {
		httpResp.Write(w, http.StatusUnauthorized, tokenNotFound)
		h.log.Error("couldn't find jwt token", sl.Op(op), sl.Err(err))
		return
	}

	userID, err := jwt.GetID(token)
	if err != nil {
		httpResp.Write(w, http.StatusUnauthorized, unauthorized)
		h.log.Error("couldn't get id", sl.Op(op), sl.Err(err))
		return
	}

	var msg = api.Message{
		ID:        0,
		PalID:     input.PalID,
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
		Message:   input.Message,
	}

	id, err := h.messager.SaveMessage(r.Context(), &msg)
	if err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't save message", sl.Op(op), sl.Err(err))
		return
	}

	httpResp.Write(w, http.StatusCreated, fmt.Sprintf("created message: id = %d", id))
}

func (h *HandlerHTTP) HandleChatRequest(w http.ResponseWriter, r *http.Request) {
	const op = scope + "HandleChatRequest"

	palIDVal, ok := r.URL.Query()[palIdKey]
	if !ok {
		httpResp.Write(w, http.StatusBadRequest, queryNotFound)
		h.log.Error("couldn't find pal_id query key", sl.Op(op))
		return
	}
	palID, err := strconv.ParseInt(palIDVal[0], 10, 64)
	if err != nil {
		httpResp.Write(w, http.StatusBadRequest, queryBadValue)
		h.log.Error("couldn't cast query val to int", sl.Op(op), sl.Err(err))
		return
	}

	offsetVal, ok := r.URL.Query()[offsetKey]
	if !ok {
		httpResp.Write(w, http.StatusBadRequest, queryNotFound)
		h.log.Error("couldn't find offset query key", sl.Op(op))
		return
	}
	offset, err := strconv.Atoi(offsetVal[0])
	if err != nil {
		httpResp.Write(w, http.StatusBadRequest, queryBadValue)
		h.log.Error("couldn't cast query val to int", sl.Op(op), sl.Err(err))
		return
	}

	limitVal, ok := r.URL.Query()[limitKey]
	if !ok {
		httpResp.Write(w, http.StatusBadRequest, queryNotFound)
		h.log.Error("couldn't find limit query key", sl.Op(op))
		return
	}
	limit, err := strconv.Atoi(limitVal[0])
	if err != nil {
		httpResp.Write(w, http.StatusBadRequest, queryBadValue)
		h.log.Error("couldn't cast query val to int", sl.Op(op), sl.Err(err))
		return
	}

	token, err := httpResp.GetToken(r)
	if err != nil {
		httpResp.Write(w, http.StatusUnauthorized, tokenNotFound)
		h.log.Error("couldn't find jwt token", sl.Op(op), sl.Err(err))
		return
	}

	userID, err := jwt.GetID(token)
	if err != nil {
		httpResp.Write(w, http.StatusUnauthorized, unauthorized)
		h.log.Error("couldn't get id from token", sl.Op(op), sl.Err(err))
		return
	}

	chat, err := h.messager.GetChat(r.Context(), userID, palID, limit, offset)
	if err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't get chat details from storage", sl.Op(op), sl.Err(err))
		return
	}

	data, err := json.Marshal(chat)
	if err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't marshal chat", sl.Op(op), sl.Err(err))
		return
	}

	if _, err = w.Write(data); err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't write chat", sl.Op(op), sl.Err(err))
		return
	}
}

func (h *HandlerHTTP) GetChats(w http.ResponseWriter, r *http.Request) {
	const op = scope + "GetChats"

	token, err := httpResp.GetToken(r)
	if err != nil {
		httpResp.Write(w, http.StatusUnauthorized, tokenNotFound)
		h.log.Error("couldn't find jwt token", sl.Op(op), sl.Err(err))
		return
	}

	userID, err := jwt.GetID(token)
	if err != nil {
		httpResp.Write(w, http.StatusUnauthorized, unauthorized)
		h.log.Error("couldn't get id from token", sl.Op(op), sl.Err(err))
		return
	}

	chats, err := h.messager.GetChats(r.Context(), userID)
	if err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't get chats", sl.Op(op), sl.Err(err))
		return
	}

	data, err := json.Marshal(chats)
	if err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't marshal chats", sl.Op(op), sl.Err(err))
		return
	}

	if _, err = w.Write(data); err != nil {
		httpResp.Write(w, http.StatusInternalServerError, internalServerError)
		h.log.Error("couldn't write chats", sl.Op(op), sl.Err(err))
		return
	}
}

func (h *HandlerHTTP) Close() error {
	const op = scope + "Close"

	if err := h.messager.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

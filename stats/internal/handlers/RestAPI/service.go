package RestAPI

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wlcmtunknwndth/messagio_test/common/httpResp"
	"github.com/wlcmtunknwndth/messagio_test/common/jwt"
	"github.com/wlcmtunknwndth/messagio_test/common/sl"
	"github.com/wlcmtunknwndth/stats/internal/domain/models"
	"log/slog"
	"net/http"
	"strconv"
)

type Storage interface {
	MessagesReceivedByUser(ctx context.Context, userID, since, to int64) (*models.UserStats, error)
	MessagesSentByUser(ctx context.Context, userID, since, to int64) (*models.UserStats, error)
	MessagesReceived(ctx context.Context, since, to int64) (*models.Stats, error)
	Close() error
}

const scope = "stats.internal.handlers.RestAPI."

type Handler struct {
	storage Storage
	log     *slog.Logger
}

func New(storage Storage, log *slog.Logger) *Handler {
	return &Handler{
		storage: storage,
		log:     log,
	}
}

const (
	sinceKey = "since"
	toKey    = "to"

	statusBadQuery            = "Bad query"
	statusUnauthorized        = "Unauthorized"
	statusInternalServerError = "Internal server error"
)

var errKeyNotFound = errors.New("key not found")

func getSinceAndTo(r *http.Request) (int64, int64, error) {
	const op = scope + "getSinceAndTo"

	sinceVals, ok := r.URL.Query()[sinceKey]
	if !ok {
		return 0, 0, fmt.Errorf("%s: %w", op, errKeyNotFound)
	}

	since, err := strconv.ParseInt(sinceVals[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	toVals, ok := r.URL.Query()[toKey]
	if !ok {
		return 0, 0, fmt.Errorf("%s: %w", op, errKeyNotFound)
	}

	to, err := strconv.ParseInt(toVals[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return since, to, nil
}

func (h *Handler) MessagesReceivedByUser(w http.ResponseWriter, r *http.Request) {
	const op = scope + "MessagesReceivedByUser"

	since, to, err := getSinceAndTo(r)
	if err != nil {
		h.log.Error("Couldn't parse query", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, statusBadQuery)
		return
	}

	token, err := httpResp.GetToken(r)
	if err != nil {
		h.log.Error("Couldn't get auth token", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	id, err := jwt.GetID(token)
	if err != nil {
		h.log.Error("Couldn't get user id", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	stats, err := h.storage.MessagesReceivedByUser(r.Context(), id, since, to)
	if err != nil {
		h.log.Error("Couldn't get statistics", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}

	msg, err := json.Marshal(stats)
	if err != nil {
		h.log.Error("Couldn't marshal stats", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}

	if _, err = w.Write(msg); err != nil {
		h.log.Error("Couldn't write respone", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}
	return
}

func (h *Handler) MessagesReceived(w http.ResponseWriter, r *http.Request) {
	const op = scope + "MessagesReceived"

	since, to, err := getSinceAndTo(r)
	if err != nil {
		h.log.Error("Couldn't parse query", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, statusBadQuery)
		return
	}

	token, err := httpResp.GetToken(r)
	if err != nil {
		h.log.Error("Couldn't get auth token", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	isAdmin, err := jwt.IsAdmin(token)
	if err != nil {
		h.log.Error("Couldn't get user id", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	if !isAdmin {
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	stats, err := h.storage.MessagesReceived(r.Context(), since, to)
	if err != nil {
		h.log.Error("Couldn't get stats", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}

	msg, err := json.Marshal(stats)
	if err != nil {
		h.log.Error("Couldn't marshal stats", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}

	if _, err = w.Write(msg); err != nil {
		h.log.Error("Couldn't write stats", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}
}

func (h *Handler) MessagesSentByUser(w http.ResponseWriter, r *http.Request) {
	const op = scope + "MessagesReceivedByUser"

	since, to, err := getSinceAndTo(r)
	if err != nil {
		h.log.Error("Couldn't parse query", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusBadRequest, statusBadQuery)
		return
	}

	token, err := httpResp.GetToken(r)
	if err != nil {
		h.log.Error("Couldn't get auth token", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	id, err := jwt.GetID(token)
	if err != nil {
		h.log.Error("Couldn't get user id", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusUnauthorized, statusUnauthorized)
		return
	}

	stats, err := h.storage.MessagesSentByUser(r.Context(), id, since, to)
	if err != nil {
		h.log.Error("Couldn't get statistics", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}

	msg, err := json.Marshal(stats)
	if err != nil {
		h.log.Error("Couldn't marshal stats", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}

	if _, err = w.Write(msg); err != nil {
		h.log.Error("Couldn't write response", sl.Op(op), sl.Err(err))
		httpResp.Write(w, http.StatusInternalServerError, statusInternalServerError)
		return
	}
	return
}

func (h *Handler) Close() error {
	const op = scope + "Close"

	if err := h.storage.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

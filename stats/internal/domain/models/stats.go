package models

type UserStats struct {
	UserID int64 `json:"user_id"`
	Stats
}

type Stats struct {
	MessagesSent   int64 `json:"messages"`
	StartTimestamp int64 `json:"start"`
	EndTimestamp   int64 `json:"finish,omitempty"`
}

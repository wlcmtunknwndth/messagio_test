package models

type UserStats struct {
	UserID int64 `json:"user_id" gorm:"unique;uniqueIndex"`
	Stats
}

type Stats struct {
	//ID           int64 `json:"id" gorm:"primaryKey;uniqueIndex;"`
	MessagesCounter int64 `json:"messages"`
	Since           int64 `json:"since"`
	To              int64 `json:"to"`
}

type MsgCount struct {
	ID        int64 `json:"id" gorm:"primaryKey;uniqueIndex;autoIncrement:true"`
	UserID    int64 `json:"user_id" gorm:"not null;index"`
	PalID     int64 `json:"pal_id" gorm:"not null;index"`
	CreatedAt int64 `json:"created_at" gorm:"not null;index"`
}

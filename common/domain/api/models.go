package api

type Message struct {
	ID        int64  `json:"id" gorm:"primaryKey;uniqueIndex;autoIncrement:true"`
	PalID     int64  `json:"pal_id" gorm:"not null;index"`
	UserID    int64  `json:"user_id" gorm:"not null;index"`
	CreatedAt int64  `json:"created_at" gorm:"not null"`
	Message   string `json:"message" gorm:"not null;type:varchar(4096)"`
}

type Input struct {
	PalID   int64  `json:"pal_id"`
	Message string `json:"message"`
}

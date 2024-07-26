package api

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	PalID   int64  `json:"pal_id"`
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type Input struct {
	PalID   int64  `json:"pal_id"`
	Message string `json:"message"`
}

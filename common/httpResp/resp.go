package httpResp

import (
	"net/http"
	"time"
)

func Write(w http.ResponseWriter, statusCode int, info string) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(info))
	if err != nil {
		return
	}
}

func WriteToken(w http.ResponseWriter, token string, duration time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:    "sso-token",
		Value:   token,
		Expires: time.Now().Add(duration),
	})
}

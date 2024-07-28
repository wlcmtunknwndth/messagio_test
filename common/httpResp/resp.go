package httpResp

import (
	"fmt"
	"net/http"
	"time"
)

const (
	jwtTokenName = "sso-token"
	scope        = "common.httpResp."
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
		Name:    jwtTokenName,
		Value:   token,
		Expires: time.Now().Add(duration),
	})
}

func GetToken(r *http.Request) (string, error) {
	const op = scope + "GetToken"
	cookie, err := r.Cookie(jwtTokenName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return cookie.Value, nil
}

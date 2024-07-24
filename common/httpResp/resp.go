package httpResp

import "net/http"

func Write(w http.ResponseWriter, statusCode int, info string) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(info))
	if err != nil {
		return
	}
}

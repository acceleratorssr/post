package web

import (
	"net/http"
)

func getUserAgent(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")

	return userAgent
}

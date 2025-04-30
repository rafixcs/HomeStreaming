package middleware

import (
	"net/http"
	"time"
)

func JsonApplicationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache headers for thumbnails
		w.Header().Set("Cache-Control", "max-age=31536000, public")
		w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(http.TimeFormat))
		next.ServeHTTP(w, r)
	})
}

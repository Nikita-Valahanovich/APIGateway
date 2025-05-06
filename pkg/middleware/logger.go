package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		requestID := r.Header.Get(RequestIDHeader)
		log.Printf("[RequestID: %s] %s %s %s", requestID, r.Method, r.URL.Path, duration)
	})
}

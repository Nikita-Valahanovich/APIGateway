package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

const RequestIDHeader = "X-Request-ID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set(RequestIDHeader, requestID)
		}
		w.Header().Set(RequestIDHeader, requestID)
		next.ServeHTTP(w, r)
	})
}

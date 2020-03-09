package middleware

import (
	"errors"
	"net/http"
	"no_homomorphism/pkg/logger"
)

func PanicMiddleware(next http.Handler, log *logger.MainLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.LogError(r.Context(), "middleware", "panicMiddleware", errors.New("panic handled"))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

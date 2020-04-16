package middleware

import (
	"fmt"
	"net/http"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
)

func PanicMiddleware(next http.Handler, log *logger.MainLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.LogError(r.Context(), "middleware", "panicMiddleware", fmt.Errorf("panic handled: %v", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

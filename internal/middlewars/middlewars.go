package middlewars

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"no_homomorphism/pkg/logger"
	"time"
)

type key int

const requestId key = 1

func AccessLogMiddleware(next http.Handler, log *logger.MainLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		rid := fmt.Sprintf("%016x", rand.Int())[:5]
		log.StartReq(*r, rid)
		start := time.Now()
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			requestId,
			rid,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
		log.EndReq(start, rid)
	})
}

func GetIdFromContext(ctx context.Context, log *logger.MainLogger) string {
	rid, ok := ctx.Value(requestId).(string)
	if !ok {
		log.LogWarning("error", "middlewars", "getIdFromContext", "no rid in context")
		return ""
	}
	return rid
}

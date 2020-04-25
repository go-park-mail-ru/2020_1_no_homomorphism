package middleware

import (
	"context"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const requestId int = 1

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{w, http.StatusOK}
}

func (srw *statusResponseWriter) WriteHeader(code int) {
	srw.statusCode = code
	srw.ResponseWriter.WriteHeader(code)
}

var hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
}, []string{"status", "path"})

func AccessLogMiddleware(next http.Handler, log *logger.MainLogger) http.Handler {
	prometheus.MustRegister(hits)

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
		srw := NewStatusResponseWriter(w)
		next.ServeHTTP(srw, r.WithContext(ctx))

		hits.WithLabelValues(strconv.Itoa(srw.statusCode), r.URL.String()).Inc()
		log.EndReq(start, ctx)
	})
}

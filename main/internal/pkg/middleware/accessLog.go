package middleware

import (
	"context"
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/config"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

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

var (
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path", "method"})

	timings = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "method_timings",
		Help: "Per method timing",
	}, []string{"method"})
)

func AccessLogMiddleware(next http.Handler, log *logger.MainLogger) http.Handler {
	prometheus.MustRegister(hits, timings)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		rid := fmt.Sprintf("%016x", rand.Int())[:5]
		log.StartReq(*r, rid)

		start := time.Now()

		ctx := r.Context()
		ctx = context.WithValue(ctx,
			config.RequestID,
			rid,
		)
		srw := NewStatusResponseWriter(w)
		next.ServeHTTP(srw, r.WithContext(ctx))

		hits.
			WithLabelValues(strconv.Itoa(srw.statusCode), r.URL.String(), r.Method).
			Inc()

		timings.
			WithLabelValues(r.URL.String()).
			Observe(time.Since(start).Seconds())

		log.EndReq(start, ctx)
	})
}

package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/pkg/logger"
	"strconv"
)

func GetBoundedVars(next http.HandlerFunc, log *logger.MainLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id, okId := vars["id"]
		startVar, okStart := vars["start"]
		endVar, okEnd := vars["end"]

		if !okId || !okStart || !okEnd {
			log.HttpInfo(r.Context(), "no data in mux vars", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		start, err1 := strconv.ParseUint(startVar, 10, 32)
		end, err2 := strconv.ParseUint(endVar, 10, 32)

		if err1 != nil || err2 != nil {
			log.HttpInfo(r.Context(), "failed to parse start or end parameters", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "id", id)
		ctx = context.WithValue(ctx, "start", start)
		ctx = context.WithValue(ctx, "end", end)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

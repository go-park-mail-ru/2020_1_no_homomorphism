package middleware

import (
	"context"
	"github.com/2020_1_no_homomorphism/no_homo_main/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type VarsKeys string

const (
	Id    VarsKeys = "id"
	Start VarsKeys = "start"
	End   VarsKeys = "end"
)

func BoundedVars(next http.HandlerFunc, log *logger.MainLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, Id, id)
		ctx = context.WithValue(ctx, Start, start)
		ctx = context.WithValue(ctx, End, end)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

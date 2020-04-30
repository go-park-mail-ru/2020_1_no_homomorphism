package middleware

import (
	"context"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

type VarsPair struct {
	Key   string
	Value string
}

func AuthMiddlewareMock(next http.HandlerFunc, auth bool, user models.User, sessionId string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "isAuth", auth)
		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "session_id", sessionId)
		ctx = context.WithValue(ctx, "isCSRFTokenCorrect", true)
		next(w, r.WithContext(ctx))
	}
}

func SetMuxVars(next http.HandlerFunc, key, value string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := mux.SetURLVars(r, map[string]string{key: value})
		next(w, req)
	}
}

func SetUnlimitedVars(next http.HandlerFunc, params ...VarsPair) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		varsMap := make(map[string]string, len(params))
		for _, elem := range params {
			varsMap[elem.Key] = elem.Value
		}
		req := mux.SetURLVars(r, varsMap)
		next(w, req)
	}
}

func SetTripleVars(next http.HandlerFunc, id, start, end string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := mux.SetURLVars(r, map[string]string{"id": id, "start": start, "end": end})
		next(w, req)
	}
}

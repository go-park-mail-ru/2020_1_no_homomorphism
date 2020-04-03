package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"no_homomorphism/internal/pkg/models"
)

func AuthMiddlewareMock(next http.HandlerFunc, auth bool, user models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "isAuth", auth)
		ctx = context.WithValue(ctx, "user", user)
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

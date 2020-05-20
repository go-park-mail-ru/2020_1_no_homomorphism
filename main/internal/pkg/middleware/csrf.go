package middleware

import (
	"context"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CsrfMiddleware struct {
	CSRF csrf.UseCase
}

type CSRFCtxKey string

const CSRFTokenCorrect CSRFCtxKey = "isCSRFTokenCorrect"

func NewCsrfMiddleware(csrf csrf.UseCase) CsrfMiddleware {
	return CsrfMiddleware{
		CSRF: csrf,
	}
}

func (m *CsrfMiddleware) CSRFCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sid, ok := ctx.Value(SessionIDKey).(string)
		if !ok {
			logrus.Error("failed to get from ctx")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		CSRFToken := r.Header.Get("Csrf-Token")
		ok, err := m.CSRF.Check(sid, CSRFToken)
		if err != nil || !ok {
			ctx = context.WithValue(ctx, CSRFTokenCorrect, false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, CSRFTokenCorrect, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

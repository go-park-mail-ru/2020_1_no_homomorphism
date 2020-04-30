package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/csrf"
)

type CsrfMiddleware struct {
	CSRF csrf.UseCase
}

func NewCsrfMiddleware(csrf csrf.UseCase) CsrfMiddleware {
	return CsrfMiddleware{
		CSRF: csrf,
	}
}

func (m *CsrfMiddleware) CSRFCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sid, ok := ctx.Value("session_id").(string)
		if !ok {
			logrus.Error("failed to get from ctx")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		CSRFToken := r.Header.Get("Csrf-Token")
		ok, err := m.CSRF.Check(sid, CSRFToken)
		if err != nil || !ok {
			ctx = context.WithValue(ctx, "isCSRFTokenCorrect", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, "isCSRFTokenCorrect", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

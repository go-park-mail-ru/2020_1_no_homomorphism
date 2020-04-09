package middleware

import (
	"context"
	"net/http"
	"no_homomorphism/internal/pkg/csrf"
)

type CsrfMiddleware struct {
	CSRF csrf.UseCase
}

func NewCsrfMiddleware(csrf csrf.UseCase) CsrfMiddleware {
	return CsrfMiddleware{
		CSRF: csrf,
	}
}

func (m *CsrfMiddleware) CSRFCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !ctx.Value("isAuth").(bool) {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sid := ctx.Value("session_id").(string)
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

package middleware

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/internal/pkg/user"
)

type Middleware struct {
	SessionUC session.UseCase
	UserUC    user.UseCase
	TrackUC   track.UseCase
}

func NewMiddleware(suc session.UseCase, uuc user.UseCase, tuc track.UseCase) *Middleware {
	return &Middleware{
		SessionUC: suc,
		UserUC:    uuc,
		TrackUC:   tuc,
	}
}

func (m *Middleware) CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sid, err := uuid.FromString(cookie.Value)
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		userInBase, err := m.SessionUC.GetUserBySessionID(sid)
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, "isAuth", true)
		ctx = context.WithValue(ctx, "user", userInBase)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

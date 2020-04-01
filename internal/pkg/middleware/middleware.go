package middleware

import (
	"context"
	"net/http"
	"no_homomorphism/internal/pkg/playlist"

	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/internal/pkg/user"
)

type Middleware struct {
	SessionDelivery session.Delivery
	UserUC          user.UseCase
	TrackUC         track.UseCase
	PlaylistUC      playlist.UseCase
}

func NewMiddleware(sd session.Delivery, uuc user.UseCase, tuc track.UseCase, puc playlist.UseCase) Middleware {
	return Middleware{
		SessionDelivery: sd,
		UserUC:          uuc,
		TrackUC:         tuc,
		PlaylistUC:      puc,
	}
}

func (m *Middleware) CheckAuthMiddleware(next http.Handler) http.Handler { //todo write logs
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		userLogin, err := m.SessionDelivery.GetLoginBySessionID(cookie.Value)
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		profile, err := m.UserUC.GetUserByLogin(userLogin)
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, "isAuth", true)
		ctx = context.WithValue(ctx, "user", profile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

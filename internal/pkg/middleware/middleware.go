package middleware

import (
	"context"
	"net/http"

	"no_homomorphism/internal/pkg/csrf"
	"no_homomorphism/internal/pkg/playlist"

	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/track"
	"no_homomorphism/internal/pkg/user"
)

type MiddlewareManager struct {
	SessionDelivery session.Delivery
	UserUC          user.UseCase
	TrackUC         track.UseCase
	PlaylistUC      playlist.UseCase
	CSRF            csrf.CryptToken
}

func NewMiddlewareManager(sd session.Delivery, uuc user.UseCase, tuc track.UseCase, puc playlist.UseCase, csrfToken csrf.CryptToken) MiddlewareManager {
	return MiddlewareManager{
		SessionDelivery: sd,
		UserUC:          uuc,
		TrackUC:         tuc,
		PlaylistUC:      puc,
		CSRF:            csrfToken,
	}
}

func (m *MiddlewareManager) CheckAuthMiddleware(next http.Handler) http.Handler { //todo write logs
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
		ctx = context.WithValue(ctx, "session_id", cookie.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareManager) CSRFCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !ctx.Value("isAuth").(bool) {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sid := ctx.Value("session_id").(string)
		CSRFToken := r.FormValue("csrf_token")
		_, err := m.CSRF.Check(sid, CSRFToken)
		if err != nil {
			ctx = context.WithValue(ctx, "isCSRFTokenCorrect", false)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, "isCSRFTokenCorrect", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

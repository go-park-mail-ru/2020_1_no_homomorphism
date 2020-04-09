package middleware

import (
	"context"
	"net/http"

	"no_homomorphism/internal/pkg/session"
	"no_homomorphism/internal/pkg/user"
)

type AuthMidleware struct {
	SessionDelivery session.Delivery
	UserUC          user.UseCase
}

func NewAuthMiddleware(sd session.Delivery, uuc user.UseCase) AuthMidleware {
	return AuthMidleware{
		SessionDelivery: sd,
		UserUC:          uuc,
	}
}

func (m *AuthMidleware) AuthMiddleware(next http.HandlerFunc) http.Handler { //todo write logs
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

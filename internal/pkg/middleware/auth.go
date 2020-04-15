package middleware

import (
	"context"
	"net/http"
	"no_homomorphism/pkg/logger"

	"no_homomorphism/configs/proto/session"
	"no_homomorphism/internal/pkg/user"
)

type AuthMidleware struct {
	SessionDelivery session.AuthCheckerClient
	UserUC          user.UseCase
	Log             *logger.MainLogger
}

func NewAuthMiddleware(sd session.AuthCheckerClient, uuc user.UseCase, logger *logger.MainLogger) AuthMidleware {
	return AuthMidleware{
		SessionDelivery: sd,
		UserUC:          uuc,
		Log:             logger,
	}
}

func (m *AuthMidleware) AuthMiddleware(next http.HandlerFunc, passNext bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			m.passNext(passNext, next, w, r, ctx)
			return
		}
		sess, err := m.SessionDelivery.Check(context.Background(), &session.SessionID{ID: cookie.Value})
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			m.passNext(passNext, next, w, r, ctx)
			return
		}
		profile, err := m.UserUC.GetUserByLogin(sess.Login)
		if err != nil {
			ctx = context.WithValue(ctx, "isAuth", false)
			m.passNext(passNext, next, w, r, ctx)
			return
		}
		ctx = context.WithValue(ctx, "isAuth", true)
		ctx = context.WithValue(ctx, "user", profile)
		ctx = context.WithValue(ctx, "session_id", cookie.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMidleware) passNext(passNext bool, next http.HandlerFunc, w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if passNext {
		next.ServeHTTP(w, r.WithContext(ctx))
	} else {
		m.Log.HttpInfo(r.Context(), "permission denied: user is not auth", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

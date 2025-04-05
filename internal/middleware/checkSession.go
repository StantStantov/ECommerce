package middleware

import (
	"Stant/ECommerce/internal/domain/models"
	"context"
	"net/http"
)

func CheckSessionMiddleware(sessions models.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("session_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			session, err := sessions.Read(sessionCookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := setUserId(r.Context(), session.UserID())
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func setUserId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, "userID", id)
}

func GetUserId(ctx context.Context) (string, bool) {
	id, ok := ctx.Value("userID").(string)
	if !ok {
		return "", false
	}
	return id, true
}

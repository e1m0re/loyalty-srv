package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"e1m0re/loyalty-srv/internal/models"
)

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				slog.Error("authentication error", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if token == nil || jwt.Validate(token, ja.ValidateOptions()...) != nil {
				slog.Error("authentication error", slog.String("error", "invalid token"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			_, claims, _ := jwtauth.FromContext(r.Context())
			userIDFromClaims := models.UserID(claims["id"].(float64))
			ctx := context.WithValue(r.Context(), models.CKUserID, userIDFromClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

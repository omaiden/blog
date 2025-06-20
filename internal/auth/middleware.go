package auth

import (
	"net/http"
	"strings"

	"blog/internal/kctx"

	"github.com/moonrhythm/httpmux"
)

const (
	HeaderAuthorization = "Authorization"
	PrefixBearer        = "Bearer "
)

func AuthMiddleware() httpmux.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get(HeaderAuthorization)
			tokenStr = strings.TrimPrefix(tokenStr, PrefixBearer)

			if tokenStr == "" {
				http.Error(w, "Unauthorized", 401)
				return
			}

			claims, err := ParseJWT(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", 401)
				return
			}

			ctx := kctx.NewUserIDContext(r.Context(), claims.UserID)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

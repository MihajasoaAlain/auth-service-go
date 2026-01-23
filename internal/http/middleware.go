package httpapi

import (
	"context"
	"net/http"
	"strings"
)

func (a AuthAPI) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "unauthorized", 401)
			return
		}

		userID, err := a.JWT.Verify(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

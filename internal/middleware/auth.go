package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nevinmanoj/hostmate/internal/auth"
)

func Authorization(jwtSecret []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var token = r.Header.Get("Authorization")
			fmt.Println("Token in middleware:", token)
			claims, err := auth.ParseToken(token, jwtSecret)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

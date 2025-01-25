package middleware

import (
	"context"
	"net/http"
	"strings"
	"vc-core/internal/utils"
)

type UserIDKey struct{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Token missing", http.StatusUnauthorized)
			return
		}

		userID, err := utils.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey{}).(string); ok {
		return userID
	}
	return ""
}
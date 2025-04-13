package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"vc-server/vc-core/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey{}, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey{}).(string); ok {
		return userID
	}
	return ""
}

func GetUserObjIDFromContext(r *http.Request) (primitive.ObjectID,error) {
	if userID, ok := r.Context().Value(UserIDKey{}).(string); ok {
		if objId, ok2 := primitive.ObjectIDFromHex(userID); ok2==nil{
			return objId,nil
		}
	}
	return primitive.ObjectID{},fmt.Errorf("invalid user id")
}
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"storex/utils"
	"strings"
)

type key string

const userIDKey key = "userID"
const roleKey key = "role"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := strings.TrimSpace(r.Header.Get("Authorization"))
			if token == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// userID, err := jwt_utils.ValidateJWT(token)

			userID, role, err := utils.ValidateJWT(token)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

// Helper to get userID from context
func GetUserID(r *http.Request) string {
	id := r.Context().Value(userIDKey)
	if idStr, ok := id.(string); ok {
		return idStr
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	role := r.Context().Value(roleKey)
	if role, ok := role.(string); ok {
		return role
	}
	return ""
}

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r)
			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: Access denied", http.StatusForbidden)
		})
	}
}

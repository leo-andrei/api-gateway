package middleware

import (
	"net/http"
)

// AuthMiddleware provides authentication for routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple JWT check (can be expanded in the future)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// In a real implementation, you'd validate the JWT token here
		// For now, we're just checking for its presence

		next.ServeHTTP(w, r)
	})
}

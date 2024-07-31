// auth.go

package middleware

import (
	"context"
	"net/http"

	"examples.com/go-webapp/internal/models"
	"examples.com/go-webapp/utils/utils"
)

// AuthMiddleware is responsible for authentication
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Perform authentication checks (validate token, session, etc.)
		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate token
		userID, err := utils.ValidateToken(authToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch user from database or service
		user, err := models.GetUserByID(userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user information to request context for use in other handlers
		ctx := context.WithValue(r.Context(), "user", user)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

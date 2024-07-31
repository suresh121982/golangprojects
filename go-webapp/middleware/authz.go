// authz.go

package middleware

import (
	"net/http"

	"examples.com/go-webapp/internal/models"
)

// AuthorizationMiddleware is responsible for authorization
func AuthorizationMiddleware(requiredPermission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve user information from request context
		user, ok := r.Context().Value("user").(*models.User)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if user has the required permission
		if !user.HasPermission(requiredPermission) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// If authorized, call the next handler
		next.ServeHTTP(w, r)
	}
}

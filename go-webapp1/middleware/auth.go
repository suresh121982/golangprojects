package middleware

import (
    "net/http"
    "strings"
    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key") // Same key as in routes.go

// Claims struct for JWT
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// Auth middleware for token validation
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

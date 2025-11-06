package auth

import (
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware that validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		isHTMLRequest := strings.Contains(r.Header.Get("Accept"), "text/html")

		// Try to get token from Authorization header first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// If no token in header and it's an HTML request, try cookie
		if tokenString == "" && isHTMLRequest {
			cookie, err := r.Cookie("auth_token")
			if err == nil && cookie != nil {
				tokenString = cookie.Value
			}
		}

		// If still no token, reject the request
		if tokenString == "" {
			if isHTMLRequest {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			if isHTMLRequest {
				// Clear invalid cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "auth_token",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Store claims in request context for use in handlers
		ctx := r.Context()
		ctx = SetUserContext(ctx, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// OptionalAuthMiddleware allows requests with or without authentication
// Useful for endpoints that work differently based on auth status
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := ValidateToken(parts[1])
				if err == nil {
					ctx := SetUserContext(r.Context(), claims)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

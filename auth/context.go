package auth

import (
	"context"
)

type contextKey string

const userContextKey contextKey = "user"

// SetUserContext stores user claims in the context
func SetUserContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, userContextKey, claims)
}

// GetUserFromContext retrieves user claims from the context
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(userContextKey).(*Claims)
	return claims, ok
}

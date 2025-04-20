package internal

import (
	"context"
)

type contextKey string

const usernameKey contextKey = "username"

// SetUsernameContext adds username to the context
func SetUsernameContext(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

// GetUsernameFromContext retrieves username from the context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
} 
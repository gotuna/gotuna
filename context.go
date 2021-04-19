package gotuna

import (
	"context"
)

type ctxKeyType string

const ctxKeyUser ctxKeyType = "user"

// ErrNoUserInContext is thrown when we cannot extract the User from the current context
var ErrNoUserInContext = constError("no user in the context")

// ContextWithUser returns a context with a User value inside
func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, user)
}

// GetUserFromContext extracts and returns the User from the context
func GetUserFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(ctxKeyUser).(User)
	if !ok {
		return nil, ErrNoUserInContext
	}
	return user, nil
}

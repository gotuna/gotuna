package gotuna

import (
	"context"
	"errors"
)

type ctxKeyType string

const ctxKeyUser ctxKeyType = "user"

// ContextWithUser returns a context with a User value inside
func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, user)
}

// GetUserFromContext extracts and returns the User from the context
func GetUserFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(ctxKeyUser).(User)
	if !ok {
		return nil, errors.New("no user in the context")
	}
	return user, nil
}

package gotuna

import (
	"context"
	"net/url"
)

type ctxKeyType string

const (
	ctxKeyParams ctxKeyType = "params"
	ctxKeyUser   ctxKeyType = "user"
)

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

// ContextWithParams returns a context with all the input parameters for
// the current request query/form/route
func ContextWithParams(ctx context.Context, vars url.Values) context.Context {
	return context.WithValue(ctx, ctxKeyParams, vars)
}

// GetParam return specific request parameter (query/form/route)
func GetParam(ctx context.Context, param string) string {
	params, ok := ctx.Value(ctxKeyParams).(url.Values)
	if !ok || len(params[param]) < 1 {
		return ""
	}

	return params[param][0]
}

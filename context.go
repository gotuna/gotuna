package gotdd

import (
	"context"
	"errors"
)

const userCtxKey = "user"

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUser(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return nil, errors.New("no user in context")
	}
	return user, nil
}

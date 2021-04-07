package gotdd

import (
	"context"
	"errors"
)

const ctxKeyUser = "user"

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, user)
}

func GetUserFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(ctxKeyUser).(User)
	if !ok {
		return nil, errors.New("no user in the context")
	}
	return user, nil
}

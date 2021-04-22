package gotuna_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestStoringAndGettingRequestParamsFromContext(t *testing.T) {
	t.Run("get url parameter", func(t *testing.T) {
		params := url.Values{
			"color":    {"red"},
			"password": {"pass123"},
		}
		ctx := gotuna.ContextWithParams(context.Background(), params)
		assert.Equal(t, "red", gotuna.GetParam(ctx, "color"))
		assert.Equal(t, "pass123", gotuna.GetParam(ctx, "password"))
		assert.Equal(t, "", gotuna.GetParam(ctx, "non-existing"))
	})
}

func TestStoringAndGettingUserFromContext(t *testing.T) {

	t.Run("empty context doesn't have a user", func(t *testing.T) {
		user, err := gotuna.GetUserFromContext(context.Background())
		assert.Error(t, err)
		assert.Equal(t, nil, user)
	})

	t.Run("store and retrieve a fake user", func(t *testing.T) {
		user := doubles.MemUser1
		ctx := gotuna.ContextWithUser(context.Background(), user)
		got, err := gotuna.GetUserFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
	})
}

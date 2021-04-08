package gotuna_test

import (
	"context"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestStoringAndGettingUserFromContext(t *testing.T) {

	t.Run("empty context doesnt have a user", func(t *testing.T) {
		user, err := gotuna.GetUserFromContext(context.Background())
		assert.Error(t, err)
		assert.Equal(t, nil, user)
	})

	t.Run("store and retrieve fake user", func(t *testing.T) {
		user := doubles.MemUser1
		ctx := gotuna.ContextWithUser(context.Background(), user)
		got, err := gotuna.GetUserFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
	})
}

package gotdd_test

import (
	"context"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestStoringAndGettingUserFromContext(t *testing.T) {

	t.Run("empty context doesnt have a user", func(t *testing.T) {
		user, err := gotdd.GetUserFromContext(context.Background())
		assert.Error(t, err)
		assert.Equal(t, nil, user)
	})

	t.Run("store and retrieve fake user", func(t *testing.T) {
		user := doubles.MemUser1
		ctx := gotdd.ContextWithUser(context.Background(), user)
		got, err := gotdd.GetUserFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
	})
}

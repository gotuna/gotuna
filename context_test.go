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
		_, err := gotdd.GetUser(context.Background())
		assert.Error(t, err)
	})

	t.Run("store and retrieve fake user", func(t *testing.T) {
		user := doubles.FakeUser1
		ctx := gotdd.WithUser(context.Background(), user)
		got, err := gotdd.GetUser(ctx)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
	})
}

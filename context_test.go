package gotdd_test

import (
	"context"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestStoringAndGettingUserFromContext(t *testing.T) {

	// empty context doesnt have a user
	_, err := gotdd.GetUser(context.Background())
	assert.Error(t, err)

	// store and retrieve fake user
	user := doubles.FakeUser1
	ctx := gotdd.WithUser(context.Background(), user)
	got, err := gotdd.GetUser(ctx)
	assert.NoError(t, err)
	assert.Equal(t, user, got)
}

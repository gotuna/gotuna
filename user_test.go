package gotdd_test

import (
	"testing"

	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestAuthenticate(t *testing.T) {

	testUser := doubles.UserStub()

	t.Run("test good authentication", func(t *testing.T) {
		user, authenticated := doubles.NewUserRepositoryStub().
			Set("email", testUser.Email).
			Set("password", "pass123").
			Authenticate()

		assert.NoError(t, authenticated)
		assert.Equal(t, user.Email, testUser.Email)
	})

	t.Run("test bad password", func(t *testing.T) {
		user, authenticated := doubles.NewUserRepositoryStub().
			Set("email", testUser.Email).
			Set("password", "badbad").
			Authenticate()

		assert.Error(t, authenticated)
		assert.Equal(t, user.Email, "")
	})

	t.Run("test non existing user", func(t *testing.T) {
		user, authenticated := doubles.NewUserRepositoryStub().
			Set("email", "nonexisting@example.com").
			Set("password", "pass123").
			Authenticate()

		assert.Error(t, authenticated)
		assert.Equal(t, user.Email, "")
	})
}

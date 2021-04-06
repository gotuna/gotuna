package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestAuthenticate(t *testing.T) {

	fakeUser := doubles.FakeUser1

	t.Run("test good authentication", func(t *testing.T) {

		form := url.Values{
			"email":    {fakeUser.Email},
			"password": {"pass123"},
		}

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		user, authenticated := doubles.NewUserRepositoryStub().
			Authenticate(response, request)

		assert.NoError(t, authenticated)
		assert.Equal(t, fakeUser.GetID(), user.GetID())
	})

	t.Run("test bad password", func(t *testing.T) {

		form := url.Values{
			"email":    {fakeUser.Email},
			"password": {"bad-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		user, authenticated := doubles.NewUserRepositoryStub().
			Authenticate(response, request)

		assert.Error(t, authenticated)
		assert.Equal(t, "", user.GetID())
	})

	t.Run("test non existing user", func(t *testing.T) {
		form := url.Values{
			"email":    {"non-existing"},
			"password": {"non-existing"},
		}

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		user, authenticated := doubles.NewUserRepositoryStub().
			Authenticate(response, request)

		assert.Error(t, authenticated)
		assert.Equal(t, "", user.GetID())
	})
}

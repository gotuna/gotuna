package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestUserRepository(t *testing.T) {

	repo := gotuna.NewInMemoryUserRepository([]gotuna.InMemoryUser{
		gotuna.InMemoryUser{
			ID:    "555",
			Email: "ted@example.com",
		},
	})

	user, err := repo.GetUserByID("555")
	assert.NoError(t, err)
	assert.Equal(t, "555", user.GetID())

	userNotFound, err := repo.GetUserByID("777")
	assert.Error(t, err)
	assert.Equal(t, "", userNotFound.GetID())
}

func TestAuthenticate(t *testing.T) {

	forms := []struct {
		name   string
		form   url.Values
		userID string
		err    error
	}{
		{
			"test good authentication",
			url.Values{
				"email":    {doubles.MemUser1.Email},
				"password": {"pass123"},
			},
			doubles.MemUser1.GetID(),
			nil,
		},
		{
			"test bad password",
			url.Values{
				"email":    {doubles.MemUser1.Email},
				"password": {"bad-pass"},
			},
			"",
			gotuna.ErrWrongPassword,
		},
		{
			"test no password",
			url.Values{
				"email": {doubles.MemUser1.Email},
			},
			"",
			gotuna.ErrRequiredField,
		},
		{
			"test no email",
			url.Values{
				"password": {"pass123"},
			},
			"",
			gotuna.ErrRequiredField,
		},
		{
			"non-existing user",
			url.Values{
				"email":    {"harry@example.com"},
				"password": {"i-do-not-exist"},
			},
			"",
			gotuna.ErrCannotFindUser,
		},
	}

	for _, tt := range forms {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.form.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			response := httptest.NewRecorder()

			userRepo := doubles.NewUserRepositoryStub()
			user, authenticated := userRepo.Authenticate(response, request)

			assert.Equal(t, tt.err, authenticated)
			assert.Equal(t, tt.userID, user.GetID())

		})
	}

}

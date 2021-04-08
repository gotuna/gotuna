package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestStoringLoggedInUserToContext(t *testing.T) {
	fakeUser := doubles.MemUser1

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(fakeUser.GetID())),
		UserRepository: doubles.NewUserRepositoryStub(),
	}

	middleware := app.StoreUserToContext()

	var userInContext gotuna.User

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInContext, _ = gotuna.GetUserFromContext(r.Context())
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, fakeUser, userInContext)
}

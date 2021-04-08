package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestStoringLoggedInUserToContext(t *testing.T) {
	fakeUser := doubles.MemUser1

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{
		Session:        gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(fakeUser.GetID())),
		UserRepository: doubles.NewUserRepositoryStub(),
	}

	middleware := app.StoreUserToContext()

	var userInContext gotdd.User

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInContext, _ = gotdd.GetUserFromContext(r.Context())
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, fakeUser, userInContext)
}

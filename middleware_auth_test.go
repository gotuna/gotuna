package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestGuestIsRedirectedToTheLoginPage(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{
		Session: gotdd.NewSession(doubles.NewGorillaSessionStoreSpy("")),
	}

	middleware := app.Authenticate("/pleaselogin")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/pleaselogin", http.StatusFound)
}

func TestLoggedInUserIsRedirectedToHome(t *testing.T) {

	fakeUser := doubles.FakeUser1

	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{
		Session: gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(fakeUser.GetID())),
	}

	middleware := app.RedirectIfAuthenticated("/dashboard")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/dashboard", http.StatusFound)
}

func TestStoringLoggedInUserToContext(t *testing.T) {
	fakeUser := doubles.FakeUser1

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{
		Session:        gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(fakeUser.GetID())),
		UserRepository: doubles.NewUserRepositoryStub(),
	}

	middleware := app.StoreUserToContext()

	var userInContext gotdd.User
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInContext, _ = gotdd.GetUser(r.Context())
	}))

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, fakeUser, userInContext)
}

package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestGuestIsRedirectedToTheLoginPage(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session: gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(""), "test"),
	}

	middleware := app.Authenticate("/pleaselogin")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/pleaselogin", http.StatusFound)
}

func TestUserIsNotRedirectedToTheLoginPage(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session: gotuna.NewSession(doubles.NewGorillaSessionStoreSpy("123"), "test"),
	}

	middleware := app.Authenticate("/pleaselogin")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestUserIsRedirectedToHome(t *testing.T) {

	fakeUser := doubles.MemUser1

	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session: gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(fakeUser.GetID()), "test"),
	}

	middleware := app.RedirectIfAuthenticated("/dashboard")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/dashboard", http.StatusFound)
}

func TestGuestIsNotRedirectedToHome(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session: gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(""), "test"),
	}

	middleware := app.RedirectIfAuthenticated("/dashboard")
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

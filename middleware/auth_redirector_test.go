package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/alcalbg/gotdd/util"
)

func TestRedirections(t *testing.T) {

	t.Run("guest is redirected to the login page", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		options := util.Options{
			Session:     session.NewSession(doubles.NewGorillaSessionStoreSpy(session.GuestSID)),
			GuestRoutes: util.GuestRoutes,
		}

		authRedirector := middleware.AuthRedirector(options)
		handler := authRedirector(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(response, request)

		assert.Redirects(t, response, "/login", http.StatusFound)
	})

	t.Run("logged in user should be redirected back from the login page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		options := util.Options{
			Session:     session.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)),
			GuestRoutes: util.GuestRoutes,
		}

		authRedirector := middleware.AuthRedirector(options)
		handler := authRedirector(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(response, request)

		assert.Redirects(t, response, "/", http.StatusFound)
	})

	t.Run("requests to static resources should skip checks early", func(t *testing.T) {
		staticFile := "/file.zip"
		request, _ := http.NewRequest(http.MethodGet, staticFile, nil)
		response := httptest.NewRecorder()

		sessionStore := doubles.NewGorillaSessionStoreSpy(session.GuestSID)
		options := util.Options{
			Session:     session.NewSession(sessionStore),
			GuestRoutes: util.GuestRoutes,
		}

		authRedirector := middleware.AuthRedirector(options)
		handler := authRedirector(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(response, request)

		assert.Equal(t, sessionStore.GetCalls, 0)
		assert.Equal(t, response.Code, http.StatusOK)
	})
}

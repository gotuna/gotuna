package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestRedirections(t *testing.T) {

	t.Run("guest is redirected to login page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		session := session.NewSession(doubles.NewSessionStoreSpy(""))

		middleware := middleware.AuthRedirector(session)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(response, request)

		assert.Redirects(t, response, "/login", http.StatusFound)
	})

	t.Run("logged in user should be redirected back from login page", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		session := session.NewSession(doubles.NewSessionStoreSpy("123"))

		middleware := middleware.AuthRedirector(session)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(response, request)

		assert.Redirects(t, response, "/", http.StatusFound)
	})

	t.Run("requests to public resources should skip checks early", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/public/file.zip", nil)
		response := httptest.NewRecorder()

		sessionStore := doubles.NewSessionStoreSpy("")
		session := session.NewSession(sessionStore)

		middleware := middleware.AuthRedirector(session)
		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(response, request)

		assert.Equal(t, sessionStore.GetCalls, 0)
		assert.Equal(t, response.Code, http.StatusOK)
	})
}

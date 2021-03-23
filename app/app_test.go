package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestRoutes(t *testing.T) {
	routes := []struct {
		userSID string
		route   string
		method  string
		status  int
	}{
		{"", "/", http.MethodGet, http.StatusFound},
		{"123", "/", http.MethodGet, http.StatusOK},
		{"", "/", http.MethodPost, http.StatusMethodNotAllowed},
		{"", "/invalid", http.MethodGet, http.StatusNotFound},
		{"", "/login", http.MethodGet, http.StatusOK},
		{"123", "/login", http.MethodGet, http.StatusFound},
		{"", "/register", http.MethodGet, http.StatusOK},
		{"123", "/register", http.MethodGet, http.StatusFound},
	}

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s", r.route), func(t *testing.T) {

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			srv := app.NewServer(
				doubles.StubLogger(),
				session.NewSession(doubles.NewSessionStoreSpy(r.userSID)),
				doubles.NewUserRepository(app.User{}),
			)

			srv.Router.ServeHTTP(response, request)

			assert.Equal(t, response.Code, r.status)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Run("test show login template", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		srv := app.NewServer(
			doubles.StubLogger(),
			session.NewSession(doubles.NewSessionStoreSpy("")),
			doubles.NewUserRepository(app.User{}),
		)

		srv.Router.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Contains(t, response.Body.String(), "Log In")
	})

	t.Run("test submit login form should save user to session", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", "john@example.com")
		data.Set("password", "pass123")

		request, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		sessionStorageSpy := doubles.NewSessionStoreSpy("")

		srv := app.NewServer(
			doubles.StubLogger(),
			session.NewSession(sessionStorageSpy),
			doubles.NewUserRepository(app.User{
				SID:          "123",
				Email:        "john@example.com",
				PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm", // pass123
			}),
		)

		srv.Router.ServeHTTP(response, request)

		assert.Redirects(t, response, "/", http.StatusFound)

		s, err := sessionStorageSpy.Get(request, "")
		assert.NoError(t, err)
		assert.Equal(t, s.Values[session.UserSIDKey], "123")
	})
}

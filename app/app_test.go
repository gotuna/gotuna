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

			app.NewServer(
				doubles.NewLoggerStub(),
				session.NewSession(doubles.NewGorillaSessionStoreSpy(r.userSID)),
				doubles.NewUserRepositoryStub(app.User{}),
			).Router.ServeHTTP(response, request)

			assert.Equal(t, response.Code, r.status)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Run("show login template", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		doubles.NewServerStub().Router.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Contains(t, response.Body.String(), `action="/login"`)
	})

	t.Run("submit login with non-existing user", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", "nonexisting@example.com")
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		doubles.NewServerStub().Router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("submit login bad password", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.UserStub().Email)
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		doubles.NewServerStub().Router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("submit successful login and go to the home page", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.UserStub().Email)
		data.Set("password", "pass123")

		// step1: after successful login, user is redirected to the home page
		request := loginRequest(data)
		response := httptest.NewRecorder()
		doubles.NewServerWithCookieStoreStub().Router.ServeHTTP(response, request)
		assert.Redirects(t, response, "/", http.StatusFound)
		gotCookies := response.Result().Cookies()

		// step2: user shoud stay on the home page
		request, _ = http.NewRequest(http.MethodGet, "/", nil)
		response = httptest.NewRecorder()
		for _, c := range gotCookies {
			request.AddCookie(c)
		}
		doubles.NewServerWithCookieStoreStub().Router.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
	})
}

func TestServingPublicFiles(t *testing.T) {
	t.Run("get robots file", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/public/robots.txt", nil)
		response := httptest.NewRecorder()

		doubles.NewServerStub().Router.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Contains(t, response.Body.String(), "robots")
	})

	t.Run("get non-existing file should return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/public/non-existing.txt", nil)
		response := httptest.NewRecorder()

		doubles.NewServerStub().Router.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func loginRequest(form url.Values) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/models"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/alcalbg/gotdd/util"
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
		t.Run(fmt.Sprintf("test route %s for %s", r.route, r.userSID), func(t *testing.T) {

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			app.NewApp(util.Options{
				Logger:         doubles.NewLoggerStub(),
				FS:             doubles.NewFileSystemStub(nil),
				Session:        session.NewSession(doubles.NewGorillaSessionStoreSpy(r.userSID)),
				UserRepository: doubles.NewUserRepositoryStub(models.User{}),
			}).ServeHTTP(response, request)

			assert.Equal(t, response.Code, r.status)
		})
	}
}

func TestServingStaticFilesFromPublicFolder(t *testing.T) {

	files := map[string][]byte{
		"somedir/image.jpg": nil,
	}

	t.Run("return valid static file from root", func(t *testing.T) {

		app := app.NewApp(util.Options{
			Logger:         doubles.NewLoggerStub(),
			FS:             doubles.NewFileSystemStub(files),
			Session:        session.NewSession(doubles.NewGorillaSessionStoreSpy(session.GuestSID)),
			UserRepository: doubles.NewUserRepositoryStub(doubles.UserStub()),
		})

		r, _ := http.NewRequest(http.MethodGet, "/somedir/image.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
	})

	t.Run("return valid static file from prefixed path", func(t *testing.T) {

		app := app.NewApp(util.Options{
			Logger:         doubles.NewLoggerStub(),
			FS:             doubles.NewFileSystemStub(files),
			Session:        session.NewSession(doubles.NewGorillaSessionStoreSpy(session.GuestSID)),
			UserRepository: doubles.NewUserRepositoryStub(doubles.UserStub()),
			StaticPrefix:   "/publicprefix",
		})

		r, _ := http.NewRequest(http.MethodGet, "/publicprefix/somedir/image.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
	})

	t.Run("return 404 on non existing file", func(t *testing.T) {

		app := app.NewApp(util.Options{
			Logger:         doubles.NewLoggerStub(),
			FS:             doubles.NewFileSystemStub(files),
			Session:        session.NewSession(doubles.NewGorillaSessionStoreSpy(session.GuestSID)),
			UserRepository: doubles.NewUserRepositoryStub(doubles.UserStub()),
		})

		r, _ := http.NewRequest(http.MethodGet, "/pic/non-existing.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNotFound)
	})

}

func TestLogin(t *testing.T) {

	htmlNeedle := `action="/login"`

	t.Run("show login template", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		doubles.NewAppStub().ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit login with non-existing user", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", "nonexisting@example.com")
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		doubles.NewAppStub().ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit login with bad password", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.UserStub().Email)
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		doubles.NewAppStub().ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit successful login and go to the home page", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.UserStub().Email)
		data.Set("password", "pass123")

		app := doubles.NewAppWithCookieStoreStub()

		// step1: after successful login, user is redirected to the home page
		request := loginRequest(data)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		assert.Redirects(t, response, "/", http.StatusFound)
		gotCookies := response.Result().Cookies()

		// step2: user shoud stay on the home page
		request, _ = http.NewRequest(http.MethodGet, "/", nil)
		response = httptest.NewRecorder()
		for _, c := range gotCookies {
			request.AddCookie(c)
		}
		app.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
	})
}

func TestLogout(t *testing.T) {

	user := doubles.UserStub()

	app := app.NewApp(util.Options{
		Logger:         doubles.NewLoggerStub(),
		FS:             doubles.NewFileSystemStub(nil),
		Session:        session.NewSession(doubles.NewGorillaSessionStoreSpy(user.SID)),
		UserRepository: doubles.NewUserRepositoryStub(user),
	})

	// first, let's make sure we're logged in
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK)

	// try to log out
	request, _ = http.NewRequest(http.MethodPost, "/logout", nil)
	response = httptest.NewRecorder()
	app.ServeHTTP(response, request)
	assert.Redirects(t, response, "/login", http.StatusFound)

	// make sure we can't reach home page anymore
	request, _ = http.NewRequest(http.MethodGet, "/", nil)
	response = httptest.NewRecorder()
	app.ServeHTTP(response, request)
	assert.Redirects(t, response, "/login", http.StatusFound)
}

func loginRequest(form url.Values) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

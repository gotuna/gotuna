package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/cmd/main"
	"github.com/alcalbg/gotdd/cmd/main/views"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/gorilla/sessions"
)

func TestRoutes(t *testing.T) {
	routes := []struct {
		userID string
		route  string
		method string
		status int
	}{
		{"", "/", http.MethodGet, http.StatusFound},
		{doubles.FakeUser1.GetID(), "/", http.MethodGet, http.StatusOK},
		{"", "/", http.MethodPost, http.StatusMethodNotAllowed},
		{"", "/invalid", http.MethodGet, http.StatusNotFound},
		{"", "/login", http.MethodGet, http.StatusOK},
		{doubles.FakeUser1.GetID(), "/login", http.MethodGet, http.StatusFound},
		{doubles.FakeUser1.GetID(), "/profile", http.MethodGet, http.StatusOK},
		{doubles.FakeUser2.GetID(), "/profile", http.MethodGet, http.StatusOK},
	}

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s for %s", r.route, r.userID), func(t *testing.T) {

			request := httptest.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			app := main.MakeApp(gotdd.App{
				Session:        gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(r.userID)),
				Static:         doubles.NewFileSystemStub(map[string][]byte{}),
				UserRepository: doubles.NewUserRepositoryStub(),
				Views:          views.EmbededViews,
			})
			app.Router.ServeHTTP(response, request)

			assert.Equal(t, r.status, response.Code)
		})
	}
}

func TestServingStaticFilesFromPublicFolder(t *testing.T) {

	files := map[string][]byte{
		"somedir/image.jpg": nil,
	}

	t.Run("return valid static file from root", func(t *testing.T) {

		app := main.MakeApp(gotdd.App{
			Static: doubles.NewFileSystemStub(files),
		})

		r := httptest.NewRequest(http.MethodGet, "/somedir/image.jpg", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("return valid static file from prefixed path", func(t *testing.T) {

		app := main.MakeApp(gotdd.App{
			Static:       doubles.NewFileSystemStub(files),
			StaticPrefix: "/publicprefix",
		})

		r := httptest.NewRequest(http.MethodGet, "/publicprefix/somedir/image.jpg", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("return 404 on non existing file", func(t *testing.T) {

		app := main.MakeApp(gotdd.App{
			Static: doubles.NewFileSystemStub(files),
		})

		r := httptest.NewRequest(http.MethodGet, "/pic/non-existing.jpg", nil)
		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

func TestLogin(t *testing.T) {

	htmlNeedle := `action="/login"`

	t.Run("show login template", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		app := main.MakeApp(gotdd.App{
			Session: gotdd.NewSession(sessions.NewCookieStore([]byte("abc"))),
			Views:   views.EmbededViews,
		})
		app.Router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit login with non-existing user", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", "nonexisting@example.com")
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		app := main.MakeApp(gotdd.App{
			Session:        gotdd.NewSession(sessions.NewCookieStore([]byte("abc"))),
			UserRepository: doubles.NewUserRepositoryStub(),
			Views:          views.EmbededViews,
		})
		app.Router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit login with bad password", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.FakeUser1.Email)
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		app := main.MakeApp(gotdd.App{
			Session:        gotdd.NewSession(sessions.NewCookieStore([]byte("abc"))),
			UserRepository: doubles.NewUserRepositoryStub(),
			Views:          views.EmbededViews,
		})
		app.Router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit successful login and go to the home page", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.FakeUser1.Email)
		data.Set("password", "pass123")

		app := main.MakeApp(gotdd.App{
			Session:        gotdd.NewSession(sessions.NewCookieStore([]byte("abc"))),
			UserRepository: doubles.NewUserRepositoryStub(),
			Views:          views.EmbededViews,
		})

		// step1: after successful login, user is redirected to the home page
		request := loginRequest(data)
		response := httptest.NewRecorder()
		app.Router.ServeHTTP(response, request)
		assert.Redirects(t, response, "/", http.StatusFound)
		gotCookies := response.Result().Cookies()

		// step2: user shoud stay on the home page
		request = httptest.NewRequest(http.MethodGet, "/", nil)
		response = httptest.NewRecorder()
		for _, c := range gotCookies {
			request.AddCookie(c)
		}
		app.Router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestLogout(t *testing.T) {

	user := doubles.FakeUser1

	app := main.MakeApp(gotdd.App{
		Session:        gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(user.GetID())),
		Views:          views.EmbededViews,
		UserRepository: doubles.NewUserRepositoryStub(),
	})

	// first, let's make sure we're logged in
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	app.Router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// try to log out
	request = httptest.NewRequest(http.MethodPost, "/logout", nil)
	response = httptest.NewRecorder()
	app.Router.ServeHTTP(response, request)
	assert.Redirects(t, response, "/login", http.StatusFound)

	// make sure we can't reach home page anymore
	request = httptest.NewRequest(http.MethodGet, "/", nil)
	response = httptest.NewRecorder()
	app.Router.ServeHTTP(response, request)
	assert.Redirects(t, response, "/login", http.StatusFound)
}

func loginRequest(form url.Values) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request
}

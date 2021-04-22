package fullapp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/examples/fullapp"
	"github.com/gotuna/gotuna/examples/fullapp/views"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestRoutes(t *testing.T) {
	routes := []struct {
		userID string
		route  string
		method string
		status int
	}{
		{"", "/", http.MethodGet, http.StatusFound},
		{doubles.MemUser1.GetID(), "/", http.MethodGet, http.StatusOK},
		{"", "/", http.MethodPost, http.StatusMethodNotAllowed},
		{"", "/invalid", http.MethodGet, http.StatusNotFound},
		{"", "/login", http.MethodGet, http.StatusOK},
		{"", "/error", http.MethodGet, http.StatusInternalServerError},
		{doubles.MemUser1.GetID(), "/setlocale/en-US", http.MethodGet, http.StatusFound},
		{doubles.MemUser1.GetID(), "/login", http.MethodGet, http.StatusFound},
		{doubles.MemUser1.GetID(), "/profile", http.MethodGet, http.StatusOK},
		{doubles.MemUser2.GetID(), "/profile", http.MethodGet, http.StatusOK},
		{doubles.MemUser1.GetID(), "/adduser", http.MethodGet, http.StatusOK},
		{doubles.MemUser1.GetID(), "/adduser", http.MethodPost, http.StatusOK},
	}

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s for %s", r.route, r.userID), func(t *testing.T) {

			request := httptest.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			app := fullapp.MakeApp(fullapp.App{
				gotuna.App{
					Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(r.userID), "test"),
					Static:         doubles.NewFileSystemStub(map[string][]byte{}),
					UserRepository: doubles.NewUserRepositoryStub(),
					ViewFiles:      views.EmbededViews,
				},
				middlewareCsrfStub,
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

	app := fullapp.MakeApp(fullapp.App{
		gotuna.App{
			Static:       doubles.NewFileSystemStub(files),
			StaticPrefix: "/publicprefix",
		},
		middlewareCsrfStub,
	})

	r := httptest.NewRequest(http.MethodGet, "/publicprefix/somedir/image.jpg", nil)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin(t *testing.T) {

	htmlNeedle := `action="/login"`

	t.Run("show login template", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/login", nil)
		response := httptest.NewRecorder()

		app := fullapp.MakeApp(fullapp.App{
			gotuna.App{
				Session:   gotuna.NewSession(sessions.NewCookieStore([]byte("abc")), "test"),
				ViewFiles: views.EmbededViews,
			},
			middlewareCsrfStub,
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
		app := fullapp.MakeApp(fullapp.App{
			gotuna.App{
				Session:        gotuna.NewSession(sessions.NewCookieStore([]byte("abc")), "test"),
				UserRepository: doubles.NewUserRepositoryStub(),
				ViewFiles:      views.EmbededViews,
			},
			middlewareCsrfStub,
		})
		app.Router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit login with bad password", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.MemUser1.Email)
		data.Set("password", "bad")

		request := loginRequest(data)
		response := httptest.NewRecorder()
		app := fullapp.MakeApp(fullapp.App{
			gotuna.App{
				Session:        gotuna.NewSession(sessions.NewCookieStore([]byte("abc")), "test"),
				UserRepository: doubles.NewUserRepositoryStub(),
				ViewFiles:      views.EmbededViews,
			},
			middlewareCsrfStub,
		})
		app.Router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Contains(t, response.Body.String(), htmlNeedle)
	})

	t.Run("submit successful login and go to the home page", func(t *testing.T) {
		data := url.Values{}
		data.Set("email", doubles.MemUser1.Email)
		data.Set("password", "pass123")

		app := fullapp.MakeApp(fullapp.App{
			gotuna.App{
				Session:        gotuna.NewSession(sessions.NewCookieStore([]byte("abc")), "test"),
				UserRepository: doubles.NewUserRepositoryStub(),
				ViewFiles:      views.EmbededViews,
			},
			middlewareCsrfStub,
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

func TestAddNewUserToTheRepo(t *testing.T) {
	app := fullapp.MakeApp(fullapp.App{
		gotuna.App{
			Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.MemUser1.ID), "test"),
			UserRepository: doubles.NewUserRepositoryStub(),
			ViewFiles:      views.EmbededViews,
		},
		middlewareCsrfStub,
	})

	form := url.Values{}
	form.Set("id", "8888")
	form.Set("name", "Mike")
	form.Set("email", "mike@example.com")
	form.Set("password", "mike9")

	request := httptest.NewRequest(http.MethodPost, "/adduser", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.Router.ServeHTTP(response, request)
	assert.Redirects(t, response, "/", http.StatusFound)

	user, err := app.UserRepository.GetUserByID("8888")
	assert.NoError(t, err)
	assert.Equal(t, "Mike", user.(gotuna.InMemoryUser).Name)
	assert.Equal(t, "mike@example.com", user.(gotuna.InMemoryUser).Email)
	assert.Equal(t, "mike9", user.(gotuna.InMemoryUser).Password)
}

func TestLogout(t *testing.T) {

	user := doubles.MemUser1

	app := fullapp.MakeApp(fullapp.App{
		gotuna.App{
			Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(user.GetID()), "test"),
			ViewFiles:      views.EmbededViews,
			UserRepository: doubles.NewUserRepositoryStub(),
		},
		middlewareCsrfStub,
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

func middlewareCsrfStub(next http.Handler) http.Handler {
	return next
}

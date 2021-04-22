package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestStoringURLParamsToContext(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/?color=red&size=xl", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{}

	middleware := app.StoreToContext()

	color := ""
	size := ""

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color = gotuna.GetParam(r.Context(), "color")
		size = gotuna.GetParam(r.Context(), "size")
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, "red", color)
	assert.Equal(t, "xl", size)
}

func TestStoringFormParamsToContext(t *testing.T) {

	form := url.Values{
		"color": {"red"},
		"size":  {"xl"},
	}
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	app := gotuna.App{}

	middleware := app.StoreToContext()

	color := ""
	size := ""

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color = gotuna.GetParam(r.Context(), "color")
		size = gotuna.GetParam(r.Context(), "size")
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, "red", color)
	assert.Equal(t, "xl", size)
}

func TestStoringRouteParamsToContext(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/car/Škoda/octavia", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Router: gotuna.NewMuxRouter(),
	}

	middleware := app.StoreToContext()

	manufacturer := ""
	model := ""

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manufacturer = gotuna.GetParam(r.Context(), "manufacturer")
		model = gotuna.GetParam(r.Context(), "model")
	})
	app.Router.Handle("/car/{manufacturer}/{model}", handler)
	app.Router.Use(middleware)
	app.Router.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, "Škoda", manufacturer)
	assert.Equal(t, "octavia", model)
}

func TestStoringLoggedInUserToContext(t *testing.T) {
	fakeUser := doubles.MemUser1

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(fakeUser.GetID()), "test"),
		UserRepository: doubles.NewUserRepositoryStub(),
	}

	middleware := app.StoreToContext()

	var userInContext gotuna.User

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInContext, _ = gotuna.GetUserFromContext(r.Context())
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, fakeUser, userInContext)
}

func TestSkipIfWeCannotFindUser(t *testing.T) {

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(""), "test"),
		UserRepository: doubles.NewUserRepositoryStub(),
	}

	middleware := app.StoreToContext()

	var userInContext gotuna.User
	var noUserErr error

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInContext, noUserErr = gotuna.GetUserFromContext(r.Context())
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Error(t, noUserErr)
	assert.Equal(t, nil, userInContext)
}

func TestErrorIfWeCannotRetreiveAuthenticatedUserFromTheRepo(t *testing.T) {

	// logged in user...
	sess := doubles.MemUser1

	// ...is not in the repo
	repo := []gotuna.InMemoryUser{
		doubles.MemUser2,
	}

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotuna.App{
		Session:        gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(sess.GetID()), "test"),
		UserRepository: gotuna.NewInMemoryUserRepository(repo),
	}

	middleware := app.StoreToContext()

	middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})).ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
}

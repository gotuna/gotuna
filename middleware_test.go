package gotdd_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestCORS(t *testing.T) {
	request, _ := http.NewRequest(http.MethodOptions, "/sample", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{}
	CORS := app.Cors()
	handler := CORS(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.HeaderMap.Get("Access-Control-Allow-Origin"), gotdd.CORSAllowedOrigin)
	assert.Equal(t, response.HeaderMap.Get("Access-Control-Allow-Methods"), gotdd.CORSAllowedMethods)
}

func TestLogging(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/sample", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}

	app := gotdd.App{

		Logger: log.New(wlog, "", 0),
	}

	logger := app.Logging()
	handler := logger(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Contains(t, wlog.String(), "GET")
	assert.Contains(t, wlog.String(), "/sample")
}

func TestRecoveringFromPanic(t *testing.T) {

	needle := "assignment to entry in nil map"

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // this code will panic with: assignment to entry in nil map
	})

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	app := gotdd.App{
		Logger: log.New(wlog, "", 0),
	}

	recoverer := app.Recoverer()
	handler := recoverer(badHandler)

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Contains(t, response.Body.String(), needle)
	assert.Contains(t, wlog.String(), needle)
}

func TestGuestIsRedirectedToTheLoginPage(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{
		Session: gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)),
	}

	authenticate := app.Authenticate("/pleaselogin")
	handler := authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/pleaselogin", http.StatusFound)
}

func TestLoggedInUserIsRedirectedToHome(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{
		Session: gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)),
	}

	redirectIfAuthenticated := app.RedirectIfAuthenticated("/dashboard")
	handler := redirectIfAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/dashboard", http.StatusFound)
}

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

func TestLoggedInUserIsRedirectedToHome(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	options := util.Options{
		Session: session.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)),
	}

	redirectIfAuthenticated := middleware.RedirectIfAuthenticated(options, "/dashboard")
	handler := redirectIfAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, "/dashboard", http.StatusFound)
}

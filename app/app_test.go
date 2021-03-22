package app_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/app"
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
	}

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s", r.route), func(t *testing.T) {

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			srv := app.NewServer(doubles.NewLogger(), doubles.NewSessionStoreSpy(request, r.userSID))

			srv.Router.ServeHTTP(response, request)

			assert.Equal(t, response.Code, r.status)
		})
	}
}

func TestGuestIsRedirectedToLoginPage(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	srv := app.NewServer(doubles.NewLogger(), doubles.NewSessionStoreSpy(request, ""))

	srv.Router.ServeHTTP(response, request)

	assert.Redirects(t, response, "/login", http.StatusFound)
}

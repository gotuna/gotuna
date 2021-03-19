package app_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/assert"
)

func TestRoutes(t *testing.T) {
	routes := []struct {
		route  string
		method string
		status int
	}{
		{"/", http.MethodGet, http.StatusOK},
		{"/", http.MethodPost, http.StatusMethodNotAllowed},
		{"/invalid", http.MethodGet, http.StatusNotFound},
		{"/login", http.MethodGet, http.StatusOK},
	}

	srv := app.NewServer(stubLogger())

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s", r.route), func(t *testing.T) {

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			srv.Router.ServeHTTP(response, request)

			assert.Equal(t, response.Code, r.status)
		})
	}
}

func stubLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

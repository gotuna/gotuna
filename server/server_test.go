package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/assert"
	"github.com/alcalbg/gotdd/server"
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

	srv := server.NewServer()

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s", r.route), func(t *testing.T) {

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			srv.Router.ServeHTTP(response, request)

			got := response.Code
			want := r.status

			assert.Equal(t, got, want)
		})
	}
}

func TestPanicWillBeRecovered(t *testing.T) {

	badHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var x map[string]int
			x["y"] = 1 // will produce nil map panic
		})
	}

	srv := server.NewServer()

	srv.Router.Handle("/bad", badHandler())

	request, _ := http.NewRequest(http.MethodGet, "/bad", nil)
	response := httptest.NewRecorder()

	srv.Router.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusInternalServerError

	assert.Equal(t, got, want)

	gotBody := response.Body.String()

	assert.Contains(t, gotBody, server.DefaultError)

}

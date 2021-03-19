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
		{"/invalid", http.MethodGet, http.StatusNotFound},
		{"/login", http.MethodGet, http.StatusOK},
	}

	for _, r := range routes {
		t.Run(fmt.Sprintf("test route %s", r.route), func(t *testing.T) {
			server := server.NewServer()

			request, _ := http.NewRequest(r.method, r.route, nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Code
			want := r.status

			assert.Equal(t, got, want)
		})
	}

}

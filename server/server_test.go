package server_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/assert"
	"github.com/alcalbg/gotdd/server"
	"github.com/alcalbg/gotdd/util"
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

	srv := server.NewServer(getTestLogger(&bytes.Buffer{}))

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

	log := &bytes.Buffer{}
	logger := getTestLogger(log)

	badHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var x map[string]int
			x["y"] = 1 // will panic with: assignment to entry in nil map
		})
	}

	srv := server.NewServer(logger)

	srv.Router.Handle("/bad", badHandler())

	request, _ := http.NewRequest(http.MethodGet, "/bad", nil)
	response := httptest.NewRecorder()

	srv.Router.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusInternalServerError

	assert.Equal(t, got, want)

	gotBody := response.Body.String()

	assert.Contains(t, gotBody, util.DefaultError)
	assert.Contains(t, log.String(), "assignment to entry in nil map")

}

func getTestLogger(w io.Writer) *log.Logger {
	return log.New(w, "test logger: ", log.Lshortfile)
}

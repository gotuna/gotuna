package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/assert"
	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/util"
)

func TestLogging(t *testing.T) {
	req, err := http.NewRequest("GET", "/sample", nil)
	assert.NoError(t, err)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	logger := log.New(wlog, "", 0)
	middleware := middleware.Logger(logger)
	handler := middleware(http.NotFoundHandler())

	handler.ServeHTTP(response, req)

	assert.Contains(t, wlog.String(), "/sample")
}

func TestRecoveringFromPanic(t *testing.T) {

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // this will panic with: assignment to entry in nil map
	})

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	logger := log.New(wlog, "", 0)
	middleware := middleware.Logger(logger)
	handler := middleware(badHandler)

	handler.ServeHTTP(response, req)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Contains(t, response.Body.String(), util.DefaultError)
	assert.Contains(t, wlog.String(), "assignment to entry in nil map")
}

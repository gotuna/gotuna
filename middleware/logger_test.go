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

func TestRecoverer(t *testing.T) {

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // will panic with: assignment to entry in nil map
	})

	req, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	logw := &bytes.Buffer{}
	logger := log.New(logw, "", 0)
	middleware := middleware.Recoverer(logger)
	handler := middleware(badHandler)

	handler.ServeHTTP(response, req)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Contains(t, response.Body.String(), util.DefaultError)
	assert.Contains(t, logw.String(), "assignment to entry in nil map")

}

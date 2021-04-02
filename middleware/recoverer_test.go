package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/util"
)

func TestRecoveringFromPanic(t *testing.T) {

	needle := "assignment to entry in nil map"

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // this code will panic with: assignment to entry in nil map
	})

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	options := util.OptionsWithDefaults(util.Options{})
	options.Logger = log.New(wlog, "", 0)

	recoverer := middleware.Recoverer(options)
	handler := recoverer(badHandler)

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Contains(t, response.Body.String(), needle)
	assert.Contains(t, wlog.String(), needle)
}

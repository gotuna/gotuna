package gotdd_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestRecoveringFromPanic(t *testing.T) {

	needle := "assignment to entry in nil map"
	destination := "/error"

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // this code will panic with: assignment to entry in nil map
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	app := gotdd.App{
		Logger: log.New(wlog, "", 0),
	}

	recoverer := app.Recoverer(destination)
	handler := recoverer(badHandler)

	handler.ServeHTTP(response, request)

	assert.Redirects(t, response, destination, http.StatusInternalServerError)
	assert.Contains(t, wlog.String(), needle)
}

package gotuna_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
)

func TestRecoveringFromPanic(t *testing.T) {

	destination := "/error"

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boo!")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	app := gotuna.App{
		Logger: log.New(wlog, "", 0),
	}

	recoverer := app.Recoverer(destination)
	recoverer(badHandler).ServeHTTP(response, request)

	assert.Redirects(t, response, destination, http.StatusInternalServerError)
	assert.Contains(t, wlog.String(), "boo!")
}

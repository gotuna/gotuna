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

func TestLogging(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/sample", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}

	app := gotdd.App{
		Logger: log.New(wlog, "", 0),
	}

	logger := app.Logging()
	logger(http.NotFoundHandler()).ServeHTTP(response, request)

	assert.Contains(t, wlog.String(), "GET")
	assert.Contains(t, wlog.String(), "/sample")
}

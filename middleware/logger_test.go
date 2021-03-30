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

func TestLogging(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/sample", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	options := util.Options{
		Logger: log.New(wlog, "", 0),
	}
	logger := middleware.Logger(options)
	handler := logger(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Contains(t, wlog.String(), "GET")
	assert.Contains(t, wlog.String(), "/sample")
}

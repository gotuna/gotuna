package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestLogging(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/sample", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	logger := log.New(wlog, "", 0)
	locale := i18n.NewLocale(i18n.En)
	middleware := middleware.Logger(logger, locale, "")
	handler := middleware(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Contains(t, wlog.String(), "GET")
	assert.Contains(t, wlog.String(), "/sample")
}

func TestRecoveringFromPanic(t *testing.T) {

	needle := "assignment to entry in nil map"

	badHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var x map[string]int
		x["y"] = 1 // this code will panic with: assignment to entry in nil map
	})

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	wlog := &bytes.Buffer{}
	logger := log.New(wlog, "", 0)
	locale := i18n.NewLocale(i18n.En)
	middleware := middleware.Logger(logger, locale, "")
	handler := middleware(badHandler)

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Contains(t, response.Body.String(), needle)
	assert.Contains(t, wlog.String(), needle)
}

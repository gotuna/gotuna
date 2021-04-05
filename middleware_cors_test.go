package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestCORS(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/sample", nil)
	response := httptest.NewRecorder()

	app := gotdd.App{}
	CORS := app.Cors()
	handler := CORS(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Equal(t, gotdd.CORSAllowedOrigin, response.HeaderMap.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, gotdd.CORSAllowedMethods, response.HeaderMap.Get("Access-Control-Allow-Methods"))
}

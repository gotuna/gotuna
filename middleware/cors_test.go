package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/middleware"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestCORS(t *testing.T) {
	request, _ := http.NewRequest(http.MethodOptions, "/sample", nil)
	response := httptest.NewRecorder()

	CORS := middleware.Cors()
	handler := CORS(http.NotFoundHandler())

	handler.ServeHTTP(response, request)

	assert.Equal(t, response.HeaderMap.Get("Access-Control-Allow-Origin"), middleware.CORSAllowedOrigin)
	assert.Equal(t, response.HeaderMap.Get("Access-Control-Allow-Methods"), middleware.CORSAllowedMethods)
}

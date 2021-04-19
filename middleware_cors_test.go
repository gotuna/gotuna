package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
)

func TestCORSHeadersAreAddedForOptionsMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/sample", nil)
	response := httptest.NewRecorder()

	middleware := gotuna.App{}.Cors()

	middleware(http.NotFoundHandler()).ServeHTTP(response, request)

	assert.Equal(t, "*", response.HeaderMap.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "OPTIONS,HEAD,GET,POST,PUT,PATCH,DELETE", response.HeaderMap.Get("Access-Control-Allow-Methods"))
}

func TestCORSHeadersAreNotAddedForGETMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/sample", nil)
	response := httptest.NewRecorder()

	middleware := gotuna.App{}.Cors()

	middleware(http.NotFoundHandler()).ServeHTTP(response, request)

	assert.Equal(t, "", response.HeaderMap.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "", response.HeaderMap.Get("Access-Control-Allow-Methods"))
}

func TestCORSHeadersAreNotAddedForPOSTMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/sample", nil)
	response := httptest.NewRecorder()

	middleware := gotuna.App{}.Cors()

	middleware(http.NotFoundHandler()).ServeHTTP(response, request)

	assert.Equal(t, "", response.HeaderMap.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "", response.HeaderMap.Get("Access-Control-Allow-Methods"))
}

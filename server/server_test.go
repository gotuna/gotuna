package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/assert"
	"github.com/alcalbg/gotdd/server"
)

func TestHome(t *testing.T) {
	server := server.NewServer()

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusOK

	assert.Equal(t, got, want)
}

func TestNotFoundOnInvalidRoute(t *testing.T) {
	server := server.NewServer()

	request, _ := http.NewRequest(http.MethodGet, "/invalid", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusNotFound

	assert.Equal(t, got, want)
}

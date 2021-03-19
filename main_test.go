package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
)

func TestHome(t *testing.T) {
	server := gotdd.NewServer()

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusOK

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

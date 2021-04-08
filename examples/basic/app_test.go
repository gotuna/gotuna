package basic_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/examples/basic"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestHome(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	app := basic.MakeApp(gotdd.App{})
	app.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Hello World!")
}

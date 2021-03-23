package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/render"
	"github.com/alcalbg/gotdd/test/assert"
)

func TestRendering(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	tmpl := render.NewTemplate("login.html")
	tmpl.Set("title", "app")
	tmpl.Render(response, request, http.StatusOK)

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Contains(t, response.Body.String(), "Log In")
}

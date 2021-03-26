package renderer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/lang"
	"github.com/alcalbg/gotdd/renderer"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestRenderingWithCustomData(t *testing.T) {

	template := `{{define "app"}}Hello, my name is {{.Data.username }}{{end}}`
	rendered := `Hello, my name is Milos`

	w := httptest.NewRecorder()

	err := getHTMLRenderer(template).
		Set("username", "Milos").
		Render(w, "view.html")

	assert.NoError(t, err)
	assert.Equal(t, w.Body.String(), rendered)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-type"), renderer.ContentTypeHTML)
}

func TestUsingTranslation(t *testing.T) {

	lang := lang.NewTranslator(map[string]string{"car": "auto"})

	template := `{{define "app"}}Hello, this is my {{.Lang.T "car" }}{{end}}`
	rendered := `Hello, this is my auto`

	w := httptest.NewRecorder()

	renderer.NewHTMLRenderer(lang).
		Mount(
			doubles.NewFileSystemStub(
				map[string][]byte{"view.html": []byte(template)})).
		Render(w, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

func TestBadTemplateShouldThrowError(t *testing.T) {

	template := `{{define "app"}} {{.Invalid.Variable}} {{end}}`

	w := httptest.NewRecorder()

	err := getHTMLRenderer(template).Render(w, "view.html")

	assert.Error(t, err)
}

func TestUsingHelperFunctions(t *testing.T) {

	template := `{{- define "app" -}} {{uppercase "hello"}} {{- end -}}`
	rendered := `HELLO`

	w := httptest.NewRecorder()

	getHTMLRenderer(template).Render(w, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

func TestLayoutWithSubContentBlock(t *testing.T) {

	const html_layout = `{{define "app"}}<div id="wrapper">{{block "sub" .}}{{end}}</div>{{end}}`
	const html_subcontent = `{{define "sub"}}<span>Subcontent</span>{{end}}`
	const html_final = `<div id="wrapper"><span>Subcontent</span></div>`

	fs := map[string][]byte{
		"layout.html":  []byte(html_layout),
		"content.html": []byte(html_subcontent),
	}

	w := httptest.NewRecorder()

	renderer.NewHTMLRenderer(nil).
		Mount(doubles.NewFileSystemStub(fs)).
		Render(w, "layout.html", "content.html")

	assert.Equal(t, w.Body.String(), html_final)
}

func getHTMLRenderer(template string) renderer.Renderer {
	rndr := renderer.NewHTMLRenderer(nil)

	// mount a fake filesystem with a single view.html file
	rndr.Mount(
		doubles.NewFileSystemStub(
			map[string][]byte{"view.html": []byte(template)}))

	return rndr
}

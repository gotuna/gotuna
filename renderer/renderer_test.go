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

func TestRenderingLayoutWithBlockContent(t *testing.T) {

	const html_layout = `{{- define "app" -}}
<!DOCTYPE html>
  <html>
	<body>
		<div class="container">{{block "content" .}}{{end}}</div>
	</body>
	</html>
{{- end -}}`

	const html_content = `{{define "content"}} {{.Data.customvar}} {{end}}`

	const html_parsed = `<!DOCTYPE html>
  <html>
	<body>
		<div class="container"> Billy </div>
	</body>
	</html>`

	fs := map[string][]byte{
		"layout.html":  []byte(html_layout),
		"content.html": []byte(html_content),
	}

	w := httptest.NewRecorder()

	htmlRenderer := renderer.NewHTMLRenderer(nil, "layout.html", "content.html")
	htmlRenderer.Mount(doubles.NewFileSystemStub(fs))
	htmlRenderer.Set("customvar", "Billy")

	err := htmlRenderer.Render(w, http.StatusOK)
	assert.NoError(t, err)
	assert.Equal(t, w.Result().Header.Get("Content-type"), renderer.ContentTypeHTML)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.String(), html_parsed)
}

func TestRenderingWithTranslation(t *testing.T) {

	lang := lang.NewTranslator(map[string]string{"car": "auto"})

	template := `{{define "app"}}Hello, this is my {{.Lang.T "car" }}{{end}}`
	rendered := `Hello, this is my auto`

	w := httptest.NewRecorder()

	htmlRenderer := renderer.NewHTMLRenderer(lang, "view.html")
	htmlRenderer.Mount(
		doubles.NewFileSystemStub(
			map[string][]byte{"view.html": []byte(template)}))

	err := htmlRenderer.Render(w, 200)
	assert.NoError(t, err)
	assert.Equal(t, w.Body.String(), rendered)
}

func TestRenderingBadTemplateShouldThrowError(t *testing.T) {

	template := `{{define "app"}} {{.Invalid.Variable}} {{end}}`

	w := httptest.NewRecorder()

	htmlRenderer := renderer.NewHTMLRenderer(nil, "view.html")
	htmlRenderer.Mount(
		doubles.NewFileSystemStub(
			map[string][]byte{"view.html": []byte(template)}))

	err := htmlRenderer.Render(w, 200)
	assert.Error(t, err)
}

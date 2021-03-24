package renderer_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/renderer"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

const testViewFile = "test.html"

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

func TestRenderingTemplates(t *testing.T) {

	files := map[string]string{}
	files[fmt.Sprintf("%s/%s", renderer.BaseDir, renderer.LayoutFile)] = html_layout
	files[fmt.Sprintf("%s/%s", renderer.BaseDir, testViewFile)] = html_content

	w := httptest.NewRecorder()

	htmlRenderer := renderer.NewHTMLRenderer(testViewFile)
	htmlRenderer.Mount(doubles.NewFileSystemStub(files))
	htmlRenderer.Set("customvar", "Billy")

	htmlRenderer.Render(w, http.StatusOK)

	assert.Equal(t, w.Result().Header.Get("Content-type"), renderer.ContentTypeHTML)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.String(), html_parsed)
}

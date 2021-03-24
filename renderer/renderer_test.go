package renderer_test

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alcalbg/gotdd/renderer"
	"github.com/alcalbg/gotdd/test/assert"
)

const testViewFile = "test.html"
const testImage = "public/image.jpg"

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
	w := httptest.NewRecorder()

	htmlRenderer := renderer.NewHTMLRenderer(testViewFile)
	htmlRenderer.Mount(newFileSystemStub())
	htmlRenderer.Set("customvar", "Billy")

	htmlRenderer.Render(w, http.StatusOK)

	assert.Equal(t, w.Result().Header.Get("Content-type"), renderer.ContentTypeHTML)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Body.String(), html_parsed)
}

func TestServingStaticFiles(t *testing.T) {
	fileServer := renderer.ServeFiles(newFileSystemStub())

	t.Run("return a valid static file", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, testImage, nil)
		w := httptest.NewRecorder()
		fileServer.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
	})

	t.Run("return 404 on non existing file", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/pic/non-existing.jpg", nil)
		w := httptest.NewRecorder()
		fileServer.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusNotFound)
	})
}

func newFileSystemStub() *filesystemStub {

	type file struct {
		name     string
		contents string
	}

	f1 := file{
		name:     fmt.Sprintf("%s/%s", renderer.BaseDir, renderer.LayoutFile),
		contents: html_layout,
	}
	f2 := file{
		name:     fmt.Sprintf("%s/%s", renderer.BaseDir, testViewFile),
		contents: html_content,
	}
	f3 := file{
		name:     testImage,
		contents: "***",
	}

	return &filesystemStub{
		files: []string{
			f1.name,
			f2.name,
			f3.name,
		},
		contents: map[string]string{
			f1.name: f1.contents,
			f2.name: f2.contents,
			f3.name: f3.contents,
		},
	}
}

type filesystemStub struct {
	files    []string
	contents map[string]string
}

func (f *filesystemStub) Glob(pattern string) ([]string, error) {
	return f.files, nil
}

func (f *filesystemStub) Open(name string) (fs.File, error) {
	tmpfile, err := ioutil.TempFile("", "fsdemo")
	if err != nil {
		log.Fatal(err)
	}

	contents, ok := f.contents[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	tmpfile.Write([]byte(contents))
	tmpfile.Seek(0, 0)

	return tmpfile, nil
}

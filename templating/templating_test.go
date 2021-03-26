package templating_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/alcalbg/gotdd/util"
)

func TestRenderingWithCustomData(t *testing.T) {

	template := `{{define "app"}}Hello, my name is {{.Data.username }}{{end}}`
	rendered := `Hello, my name is Milos`

	w := httptest.NewRecorder()

	err := doubles.NewStubTemplatingEngine(template).
		Set("username", "Milos").
		Render(w, "view.html")

	assert.NoError(t, err)
	assert.Equal(t, w.Body.String(), rendered)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-type"), util.ContentTypeHTML)
}

func TestUsingTranslation(t *testing.T) {

	lang := i18n.NewTranslator(map[string]string{"car": "auto"})

	template := `{{define "app"}}Hello, this is my {{lang "car"}}{{end}}`
	rendered := `Hello, this is my auto`

	w := httptest.NewRecorder()

	templating.GetNativeTemplatingEngine(lang).
		Mount(
			doubles.NewFileSystemStub(
				map[string][]byte{"view.html": []byte(template)})).
		Render(w, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

func TestBadTemplateShouldThrowError(t *testing.T) {

	template := `{{define "app"}} {{.Invalid.Variable}} {{end}}`

	w := httptest.NewRecorder()

	err := doubles.NewStubTemplatingEngine(template).Render(w, "view.html")

	assert.Error(t, err)
}

func TestUsingHelperFunctions(t *testing.T) {

	template := `{{- define "app" -}} {{uppercase "hello"}} {{- end -}}`
	rendered := `HELLO`

	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(template).Render(w, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

func TestLayoutWithSubContentBlock(t *testing.T) {

	const htmlLayout = `{{define "app"}}<div id="wrapper">{{block "sub" .}}{{end}}</div>{{end}}`
	const htmlSubcontent = `{{define "sub"}}<span>Subcontent</span>{{end}}`
	const htmlFinal = `<div id="wrapper"><span>Subcontent</span></div>`

	fs := map[string][]byte{
		"layout.html":  []byte(htmlLayout),
		"content.html": []byte(htmlSubcontent),
	}

	w := httptest.NewRecorder()

	templating.GetNativeTemplatingEngine(i18n.NewTranslator(nil)).
		Mount(doubles.NewFileSystemStub(fs)).
		Render(w, "layout.html", "content.html")

	assert.Equal(t, w.Body.String(), htmlFinal)
}

func TestCanChangeContentType(t *testing.T) {

	template := `{{- define "app" -}} {{uppercase "hello"}} {{- end -}}`
	rendered := `HELLO`

	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(template).Render(w, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

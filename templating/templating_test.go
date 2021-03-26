package templating_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(template).
		Set("username", "Milos").
		Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-type"), util.ContentTypeHTML)
}

func TestUsingTranslation(t *testing.T) {

	lang := i18n.NewTranslator(map[string]string{"car": "auto"})

	template := `{{define "app"}}Hello, this is my {{lang "car"}}{{end}}`
	rendered := `Hello, this is my auto`

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	templating.GetNativeTemplatingEngine(lang).
		Mount(
			doubles.NewFileSystemStub(
				map[string][]byte{"view.html": []byte(template)})).
		Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

//func TestBadTemplateShouldPanic(t *testing.T) {
//
//	template := `{{define "app"}} {{.Invalid.Variable}} {{end}}`
//
//	r, _ := http.NewRequest(http.MethodGet, "/", nil)
//	w := httptest.NewRecorder()
//
//	doubles.NewStubTemplatingEngine(template).Render(w, r, "view.html")
//}

func TestUsingHelperFunctions(t *testing.T) {

	template := `{{- define "app" -}} {{uppercase "hello"}} {{- end -}}`
	rendered := `HELLO`

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(template).Render(w, r, "view.html")

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

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	templating.GetNativeTemplatingEngine(i18n.NewTranslator(nil)).
		Mount(doubles.NewFileSystemStub(fs)).
		Render(w, r, "layout.html", "content.html")

	assert.Equal(t, w.Body.String(), htmlFinal)
}

func TestCurrentRequestCanBeUsedInTemplates(t *testing.T) {
	form := url.Values{
		"email": {"user@example.com"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/test", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusConflict)

	tmpl := `{{define "app"}}Hello {{.Request.FormValue "email"}}{{end}}`
	want := `Hello user@example.com`

	doubles.NewStubTemplatingEngine(tmpl).Render(w, r, "view.html")
	assert.Equal(t, w.Body.String(), want)
}

func TestErrorsCanBeAdded(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	tmpl := `{{define "app"}}{{index .Errors "error1"}} / {{index .Errors "error2"}}{{end}}`
	want := `some error / other error`

	engine := doubles.NewStubTemplatingEngine(tmpl).
		AddError("error1", "some error").
		AddError("error2", "other error")
	engine.Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), want)
	assert.Equal(t, len(engine.GetErrors()), 2)
}

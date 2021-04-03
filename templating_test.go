package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestRenderingWithCustomData(t *testing.T) {

	template := `{{define "app"}}Hello, my name is {{.Data.username }}{{end}}`
	rendered := `Hello, my name is Milos`

	r := &http.Request{}
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(template, gotdd.OptionsWithDefaults(gotdd.Options{})).
		Set("username", "Milos").
		Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Result().Header.Get("Content-type"), gotdd.ContentTypeHTML)
}

func TestUsingTranslation(t *testing.T) {

	options := gotdd.Options{
		Locale: gotdd.NewLocale(map[string]map[string]string{
			"car": {
				"en-US": "auto",
			},
		}),
	}

	template := `{{define "app"}}Hello, this is my {{t "car"}}{{end}}`
	rendered := `Hello, this is my auto`

	r := &http.Request{}
	w := httptest.NewRecorder()

	gotdd.GetEngine(options).
		MountFS(
			doubles.NewFileSystemStub(
				map[string][]byte{"view.html": []byte(template)})).
		Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), rendered)
}

func TestBadTemplateShouldPanic(t *testing.T) {

	template := `{{define "app"}} {{.Invalid.Variable}} {{end}}`

	r := &http.Request{}
	w := httptest.NewRecorder()

	defer func() {
		recover()
	}()

	doubles.NewStubTemplatingEngine(template, gotdd.OptionsWithDefaults(gotdd.Options{})).
		Render(w, r, "view.html")

	t.Errorf("templating engine should panic")
}

func TestUsingHelperFunctions(t *testing.T) {

	template := `{{- define "app" -}} {{uppercase "hello"}} {{- end -}}`
	rendered := `HELLO`

	r := &http.Request{}
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(template, gotdd.OptionsWithDefaults(gotdd.Options{})).
		Render(w, r, "view.html")

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

	r := &http.Request{}
	w := httptest.NewRecorder()

	gotdd.GetEngine(gotdd.OptionsWithDefaults(gotdd.Options{})).
		MountFS(doubles.NewFileSystemStub(fs)).
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

	template := `{{define "app"}}Hello {{.Request.FormValue "email"}}{{end}}`
	want := `Hello user@example.com`

	doubles.NewStubTemplatingEngine(template, gotdd.OptionsWithDefaults(gotdd.Options{})).
		Render(w, r, "view.html")
	assert.Equal(t, w.Body.String(), want)
}

func TestErrorsCanBeAdded(t *testing.T) {
	r := &http.Request{}
	w := httptest.NewRecorder()

	template := `{{define "app"}}{{index .Errors "error1"}} / {{index .Errors "error2"}}{{end}}`
	want := `some error / other error`

	engine := doubles.NewStubTemplatingEngine(template, gotdd.OptionsWithDefaults(gotdd.Options{})).
		SetError("error1", "some error").
		SetError("error2", "other error")
	engine.Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), want)
	assert.Equal(t, len(engine.GetErrors()), 2)
}

func TestFlashMessagesAreIncluded(t *testing.T) {
	r := &http.Request{}
	w := httptest.NewRecorder()

	ses := gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID))
	ses.Flash(w, r, gotdd.FlashMessage{Kind: "success", Message: "flash one"})
	ses.Flash(w, r, gotdd.FlashMessage{Kind: "success", Message: "flash two", AutoClose: true})

	template := `{{define "app"}}{{range $el := .Flashes}} * {{$el.Message}}{{end}}{{end}}`
	want := ` * flash one * flash two`

	options := gotdd.OptionsWithDefaults(gotdd.Options{
		Session: ses,
	})

	doubles.NewStubTemplatingEngine(template, options).
		Render(w, r, "view.html")

	assert.Equal(t, w.Body.String(), want)
}

func TestIsGuestIsEvaluated(t *testing.T) {
	// TODO
}

package gotdd_test

import (
	"context"
	"html/template"
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

	tmpl := `{{define "app"}}Hello, my name is {{.Data.username }}{{end}}`
	rendered := `Hello, my name is John`

	r := &http.Request{}
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(tmpl).
		Set("username", "John").
		Render(w, r, "view.html")

	assert.Equal(t, rendered, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, gotdd.ContentTypeHTML, w.Result().Header.Get("Content-type"))
}

func TestUsingTranslation(t *testing.T) {

	tmpl := `{{define "app"}}Hello, this is my {{t "car"}}{{end}}`
	rendered := `Hello, this is my auto`

	app := gotdd.App{
		Session: gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.FakeUser1.GetID())),
		Locale: gotdd.NewLocale(map[string]map[string]string{
			"car": {
				"en-US": "auto",
			},
		}),
		ViewFiles: doubles.NewFileSystemStub(
			map[string][]byte{
				"view.html": []byte(tmpl),
			}),
	}

	r := &http.Request{}
	w := httptest.NewRecorder()

	app.Session.SetUserLocale(w, r, "en-US")

	app.NewTemplatingEngine().Render(w, r, "view.html")

	assert.Equal(t, rendered, w.Body.String())
}

func TestBadTemplateShouldPanic(t *testing.T) {

	tmpl := `{{define "app"}} {{.Invalid.Variable}} {{end}}`

	r := &http.Request{}
	w := httptest.NewRecorder()

	defer func() {
		recover()
	}()

	doubles.NewStubTemplatingEngine(tmpl).
		Render(w, r, "view.html")

	t.Errorf("templating engine should panic")
}

func TestUsingHelperFunctions(t *testing.T) {

	t.Run("test native helper function", func(t *testing.T) {
		tmpl := `{{- define "app" -}} {{static "file"}} {{- end -}}`
		rendered := `file?123`

		r := &http.Request{}
		w := httptest.NewRecorder()

		doubles.NewStubTemplatingEngine(tmpl).
			Render(w, r, "view.html")

		assert.Equal(t, rendered, w.Body.String())
	})

	t.Run("test when custom helper function is added", func(t *testing.T) {
		tmpl := `{{- define "app" -}} {{customUppercase "name"}} {{- end -}}`
		rendered := `NAME`

		r := &http.Request{}
		w := httptest.NewRecorder()

		gotdd.App{
			ViewFiles: doubles.NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(tmpl),
				}),
			ViewHelpers: template.FuncMap{
				"customUppercase": func(s string) string {
					return strings.ToUpper(s)
				},
			}}.NewTemplatingEngine().
			Render(w, r, "view.html")

		assert.Equal(t, rendered, w.Body.String())
	})
}

func TestLayoutWithSubContentBlock(t *testing.T) {

	const htmlLayout = `{{define "app"}}<div id="wrapper">{{block "sub" .}}{{end}}</div>{{end}}`
	const htmlSubcontent = `{{define "sub"}}<span>Subcontent</span>{{end}}`
	const htmlFinal = `<div id="wrapper"><span>Subcontent</span></div>`

	r := &http.Request{}
	w := httptest.NewRecorder()

	gotdd.App{
		ViewFiles: doubles.NewFileSystemStub(
			map[string][]byte{
				"layout.html":  []byte(htmlLayout),
				"content.html": []byte(htmlSubcontent),
			}),
	}.
		NewTemplatingEngine().
		Render(w, r, "layout.html", "content.html")

	assert.Equal(t, htmlFinal, w.Body.String())
}

func TestCurrentRequestCanBeUsedInTemplates(t *testing.T) {
	form := url.Values{
		"email": {"user@example.com"},
	}

	r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusConflict)

	tmpl := `{{define "app"}}Hello {{request.FormValue "email"}}{{end}}`
	want := `Hello user@example.com`

	doubles.NewStubTemplatingEngine(tmpl).
		Render(w, r, "view.html")

	assert.Equal(t, want, w.Body.String())
}

func TestErrorsCanBeAdded(t *testing.T) {
	r := &http.Request{}
	w := httptest.NewRecorder()

	tmpl := `{{define "app"}}{{index .Errors "error1"}} / {{index .Errors "error2"}}{{end}}`
	want := `some error / other error`

	engine := doubles.NewStubTemplatingEngine(tmpl).
		SetError("error1", "some error").
		SetError("error2", "other error")
	engine.Render(w, r, "view.html")

	assert.Equal(t, want, w.Body.String())
	assert.Equal(t, 2, len(engine.GetErrors()))
}

func TestFlashMessagesAreIncluded(t *testing.T) {
	r := &http.Request{}
	w := httptest.NewRecorder()

	ses := gotdd.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.FakeUser1.GetID()))
	ses.Flash(w, r, gotdd.FlashMessage{Kind: "success", Message: "flash one"})
	ses.Flash(w, r, gotdd.FlashMessage{Kind: "success", Message: "flash two", AutoClose: true})

	tmpl := `{{define "app"}}{{range $el := .Flashes}} * {{$el.Message}}{{end}}{{end}}`
	want := ` * flash one * flash two`

	gotdd.App{
		Session: ses,
		ViewFiles: doubles.NewFileSystemStub(
			map[string][]byte{
				"view.html": []byte(tmpl),
			}),
	}.
		NewTemplatingEngine().
		Render(w, r, "view.html")

	assert.Equal(t, want, w.Body.String())

	messages := ses.Flashes(w, r)
	assert.Equal(t, 0, len(messages))
}

func TestCurrentUserIsIncluded(t *testing.T) {

	tmpl := `{{define "app"}}Welcome {{currentUser.Name}}{{end}}`
	rendered := `Welcome John`

	fakeUser := doubles.FakeUser1

	ctx := gotdd.ContextWithUser(context.Background(), fakeUser)
	r := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(tmpl).
		Set("username", "John").
		Render(w, r, "view.html")

	assert.Equal(t, rendered, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, gotdd.ContentTypeHTML, w.Result().Header.Get("Content-type"))
}

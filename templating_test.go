package gotuna_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
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
	assert.Equal(t, gotuna.ContentTypeHTML, w.Result().Header.Get("Content-type"))
}

func TestUsingTranslation(t *testing.T) {

	tmpl := `{{define "app"}}Hello, this is my {{t "car"}}. {{tp "oranges" 5}}{{end}}`
	rendered := `Hello, this is my auto. There are many oranges`

	app := gotuna.App{
		Session: gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.MemUser1.GetID()), "test"),
		Locale: gotuna.NewLocale(map[string]map[string]string{
			"car": {
				"en-US": "auto",
			},
			"oranges": {
				"en-US": "There is one orange|There are many oranges",
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
		tmpl := `{{- define "app" -}} {{uppercase "name"}} {{- end -}}`
		rendered := `NAME`

		r := &http.Request{}
		w := httptest.NewRecorder()

		customHelper := func(w http.ResponseWriter, r *http.Request) (string, interface{}) {
			return "uppercase", func(s string) string {
				return strings.ToUpper(s)
			}
		}

		gotuna.App{
			ViewFiles: doubles.NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(tmpl),
				}),
			ViewHelpers: []gotuna.ViewHelperFunc{customHelper},
		}.NewTemplatingEngine().
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

	gotuna.App{
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

	ses := gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.MemUser1.GetID()), "test")
	ses.SetUserLocale(w, r, "pt-PT")
	ses.Flash(w, r, gotuna.FlashMessage{Kind: "success", Message: "flash one"})
	ses.Flash(w, r, gotuna.FlashMessage{Kind: "success", Message: "flash two", AutoClose: true})

	tmpl := `{{define "app"}}{{range $el := .Flashes}} * {{$el.Message}}{{end}}{{end}}`
	want := ` * flash one * flash two`

	gotuna.App{
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

func TestLocaleIsIncludedForLoggedInUsers(t *testing.T) {
	r := &http.Request{}
	w := httptest.NewRecorder()

	ses := gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(doubles.MemUser1.GetID()), "test")
	ses.SetUserLocale(w, r, "pt-PT")

	tmpl := `{{define "app"}}{{if not isGuest}}Hi user, your locale is {{currentLocale}}{{end}}{{end}}`
	want := `Hi user, your locale is pt-PT`

	gotuna.App{
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

	fakeUser := doubles.MemUser1

	ctx := gotuna.ContextWithUser(context.Background(), fakeUser)
	r := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	doubles.NewStubTemplatingEngine(tmpl).
		Set("username", "John").
		Render(w, r, "view.html")

	assert.Equal(t, rendered, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, gotuna.ContentTypeHTML, w.Result().Header.Get("Content-type"))
}

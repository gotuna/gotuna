package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestServingStaticFilesFromPublicFolder(t *testing.T) {

	files := map[string][]byte{
		"somedir/image.jpg": nil,
		"badfile.txt":       nil,
	}

	app := gotuna.App{
		Static: doubles.NewFileSystemStub(files),
	}

	t.Run("return valid static file from root", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "/somedir/image.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeFiles(http.HandlerFunc(http.NotFound)).ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("return 404 on non existing file", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "/pic/non-existing.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeFiles(http.HandlerFunc(http.NotFound)).ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("return 404 on bad file", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "/badfile.txt", nil)
		w := httptest.NewRecorder()
		app.ServeFiles(http.HandlerFunc(http.NotFound)).ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

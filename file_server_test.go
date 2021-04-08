package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestServingStaticFilesFromPublicFolder(t *testing.T) {

	files := map[string][]byte{
		"somedir/image.jpg": nil,
	}

	t.Run("return valid static file from root", func(t *testing.T) {

		app := gotdd.App{
			Static: doubles.NewFileSystemStub(files),
		}

		r := httptest.NewRequest(http.MethodGet, "/somedir/image.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeFiles().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("return 404 on non existing file", func(t *testing.T) {

		app := gotdd.App{
			Static: doubles.NewFileSystemStub(files),
		}

		r := httptest.NewRequest(http.MethodGet, "/pic/non-existing.jpg", nil)
		w := httptest.NewRecorder()
		app.ServeFiles().ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

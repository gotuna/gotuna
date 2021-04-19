package gotuna_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
)

func TestCreateMuxRouter(t *testing.T) {
	assert.Equal(t, mux.NewRouter().Get(""), gotuna.NewMuxRouter().Get(""))
}

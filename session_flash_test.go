package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestFlashMessages(t *testing.T) {

	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.MemUser1.GetID())
	ses := gotuna.NewSession(sessionStoreSpy, "test")
	r := &http.Request{}
	w := httptest.NewRecorder()

	// request1: no flash messages
	messages := ses.Flashes(w, r)
	assert.Equal(t, 0, len(messages))

	// request2: add two flash messages
	ses.Flash(w, r, gotuna.NewFlash("flash message one"))
	ses.Flash(w, r, gotuna.FlashMessage{Message: "flash message two", Kind: "active", AutoClose: true})

	// request3: pop flash messages
	messages = ses.Flashes(w, r)
	assert.Equal(t, 2, len(messages))
	assert.Equal(t, "flash message one", messages[0].Message)
	assert.Equal(t, "success", messages[0].Kind)
	assert.Equal(t, true, messages[0].AutoClose)
	assert.Equal(t, "active", messages[1].Kind)
	assert.Equal(t, true, messages[1].AutoClose)

	// request4: no flash messages
	messages = ses.Flashes(w, r)
	assert.Equal(t, 0, len(messages))
}

func TestIvalidFlashMessageInTheSession(t *testing.T) {
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.MemUser1.GetID())
	ses := gotuna.NewSession(sessionStoreSpy, "test")
	r := &http.Request{}
	w := httptest.NewRecorder()

	ses.Put(w, r, "_flash", "alien saved to the flash key")
	err := ses.Flash(w, r, gotuna.NewFlash("flash message one"))
	assert.Error(t, err)
}

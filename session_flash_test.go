package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestFlashMessages(t *testing.T) {

	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.FakeUser1.GetID())
	ses := gotdd.NewSession(sessionStoreSpy)
	r := &http.Request{}
	w := httptest.NewRecorder()

	// request1: no flash messages
	messages := ses.Flashes(w, r)
	assert.Equal(t, 0, len(messages))

	// request2: add two flash messages
	messages = ses.Flashes(w, r)
	ses.Flash(w, r, gotdd.NewFlash("flash message one"))
	ses.Flash(w, r, gotdd.FlashMessage{Message: "flash message two", Kind: "active", AutoClose: true})

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

package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestStoreAndRetrieveData(t *testing.T) {

	t.Run("test storing and retrieving a simple string", func(t *testing.T) {
		r := &http.Request{}
		w := httptest.NewRecorder()
		sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
		ses := gotdd.NewSession(sessionStoreSpy)

		err := ses.Put(w, r, "test", "somevalue")
		assert.NoError(t, err)

		value, err := ses.Get(w, r, "test")
		assert.NoError(t, err)
		assert.Equal(t, value, "somevalue")
	})

	t.Run("test retrieving unsaved data", func(t *testing.T) {
		r := &http.Request{}
		w := httptest.NewRecorder()
		sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
		ses := gotdd.NewSession(sessionStoreSpy)

		value, err := ses.Get(w, r, "test")
		assert.Error(t, err)
		assert.Equal(t, value, "")
	})
}

func TestReadingUserSIDFromEmptyStore(t *testing.T) {

	r := &http.Request{}
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
	ses := gotdd.NewSession(sessionStoreSpy)

	assert.Equal(t, ses.IsGuest(r), true)

	sid, err := ses.GetUserSID(r)
	assert.Error(t, err)
	assert.Equal(t, sid, gotdd.GuestSID)

	assert.Equal(t, sessionStoreSpy.SaveCalls, 0)
}

func TestSaveUserSIDAndRetrieve(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
	ses := gotdd.NewSession(sessionStoreSpy)

	err := ses.SetUserSID(w, r, doubles.UserStub().SID)
	assert.NoError(t, err)

	sid, err := ses.GetUserSID(r)
	assert.NoError(t, err)
	assert.Equal(t, sid, doubles.UserStub().SID)
	assert.Equal(t, ses.IsGuest(r), false)
	assert.Equal(t, sessionStoreSpy.SaveCalls, 1)
}

func TestDestroyActiveSession(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)
	ses := gotdd.NewSession(sessionStoreSpy)

	sid, _ := ses.GetUserSID(r)
	assert.Equal(t, sid, doubles.UserStub().SID)

	ses.DestroySession(w, r)

	sid, err := ses.GetUserSID(r)
	assert.Error(t, err)
	assert.Equal(t, sid, gotdd.GuestSID)
	assert.Equal(t, sessionStoreSpy.SaveCalls, 1)
	assert.Equal(t, sessionStoreSpy.Session.Options.MaxAge, -1)
}

func TestFlashMessages(t *testing.T) {

	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)
	ses := gotdd.NewSession(sessionStoreSpy)
	r := &http.Request{}
	w := httptest.NewRecorder()

	// request1: no flash messages
	messages, err := ses.Flashes(w, r)
	assert.NoError(t, err)
	assert.Equal(t, len(messages), 0)

	// request2: add flash messages
	messages, err = ses.Flashes(w, r)
	ses.Flash(w, r, gotdd.NewFlash("flash message one"))
	ses.Flash(w, r, gotdd.FlashMessage{Message: "flash message two", Kind: "active", AutoClose: true})

	// request3: pop flash messages
	messages, err = ses.Flashes(w, r)
	assert.NoError(t, err)
	assert.Equal(t, len(messages), 2)
	assert.Equal(t, messages[0].Message, "flash message one")
	assert.Equal(t, messages[0].Kind, "success")
	assert.Equal(t, messages[0].AutoClose, true)
	assert.Equal(t, messages[1].Kind, "active")
	assert.Equal(t, messages[1].AutoClose, true)

	// request4: no flash messages
	messages, err = ses.Flashes(w, r)
	assert.NoError(t, err)
	assert.Equal(t, len(messages), 0)
}

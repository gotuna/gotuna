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

	t.Run("test storing, retrieving, and deleting a simple string", func(t *testing.T) {
		r := &http.Request{}
		w := httptest.NewRecorder()
		sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
		ses := gotdd.NewSession(sessionStoreSpy)

		err := ses.Put(w, r, "test", "somevalue")
		assert.NoError(t, err)

		value, err := ses.Get(r, "test")
		assert.NoError(t, err)
		assert.Equal(t, "somevalue", value)

		err = ses.Delete(w, r, "test")
		assert.NoError(t, err)

		value, err = ses.Get(r, "test")
		assert.Error(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("test retrieving unsaved data", func(t *testing.T) {
		r := &http.Request{}
		sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
		ses := gotdd.NewSession(sessionStoreSpy)

		value, err := ses.Get(r, "test")
		assert.Error(t, err)
		assert.Equal(t, "", value)
	})

}

func TestReadingUserSIDFromEmptyStore(t *testing.T) {

	r := &http.Request{}
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
	ses := gotdd.NewSession(sessionStoreSpy)

	assert.Equal(t, true, ses.IsGuest(r))

	sid, err := ses.GetUserSID(r)
	assert.Error(t, err)
	assert.Equal(t, gotdd.GuestSID, sid)

	assert.Equal(t, 0, sessionStoreSpy.SaveCalls)
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
	assert.Equal(t, doubles.UserStub().SID, sid)
	assert.Equal(t, false, ses.IsGuest(r))
	assert.Equal(t, 1, sessionStoreSpy.SaveCalls)
}

func TestDestroyActiveSession(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)
	ses := gotdd.NewSession(sessionStoreSpy)

	sid, _ := ses.GetUserSID(r)
	assert.Equal(t, doubles.UserStub().SID, sid)

	ses.DestroySession(w, r)

	sid, err := ses.GetUserSID(r)
	assert.Error(t, err)
	assert.Equal(t, gotdd.GuestSID, sid)
	assert.Equal(t, 1, sessionStoreSpy.SaveCalls)
	assert.Equal(t, -1, sessionStoreSpy.Session.Options.MaxAge)
}

func TestFlashMessages(t *testing.T) {

	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)
	ses := gotdd.NewSession(sessionStoreSpy)
	r := &http.Request{}
	w := httptest.NewRecorder()

	// request1: no flash messages
	messages, err := ses.Flashes(w, r)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(messages))

	// request2: add flash messages
	messages, err = ses.Flashes(w, r)
	ses.Flash(w, r, gotdd.NewFlash("flash message one"))
	ses.Flash(w, r, gotdd.FlashMessage{Message: "flash message two", Kind: "active", AutoClose: true})

	// request3: pop flash messages
	messages, err = ses.Flashes(w, r)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(messages))
	assert.Equal(t, "flash message one", messages[0].Message)
	assert.Equal(t, "success", messages[0].Kind)
	assert.Equal(t, true, messages[0].AutoClose)
	assert.Equal(t, "active", messages[1].Kind)
	assert.Equal(t, true, messages[1].AutoClose)

	// request4: no flash messages
	messages, err = ses.Flashes(w, r)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(messages))
}

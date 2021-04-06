package gotdd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

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

func TestSettingUserLocale(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(gotdd.GuestSID)
	ses := gotdd.NewSession(sessionStoreSpy)

	assert.Equal(t, "", ses.GetUserLocale(r))

	err := ses.SetUserLocale(w, r, "fr-FR")
	assert.NoError(t, err)

	assert.Equal(t, "fr-FR", ses.GetUserLocale(r))
	assert.Equal(t, 1, sessionStoreSpy.SaveCalls)
}

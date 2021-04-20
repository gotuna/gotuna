package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestReadingUserIDFromEmptyStore(t *testing.T) {

	r := &http.Request{}
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy("")
	ses := gotuna.NewSession(sessionStoreSpy, "test")

	assert.Equal(t, true, ses.IsGuest(r))

	id, err := ses.GetUserID(r)
	assert.Error(t, err)
	assert.Equal(t, "", id)

	assert.Equal(t, 0, sessionStoreSpy.SaveCalls)
}

func TestSaveUserIDAndRetrieve(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy("")
	ses := gotuna.NewSession(sessionStoreSpy, "test")

	err := ses.SetUserID(w, r, doubles.MemUser1.GetID())
	assert.NoError(t, err)

	id, err := ses.GetUserID(r)
	assert.NoError(t, err)
	assert.Equal(t, doubles.MemUser1.GetID(), id)
	assert.Equal(t, false, ses.IsGuest(r))
	assert.Equal(t, 1, sessionStoreSpy.SaveCalls)
}

func TestSettingUserLocale(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy("")
	ses := gotuna.NewSession(sessionStoreSpy, "test")

	assert.Equal(t, "", ses.GetUserLocale(r))

	err := ses.SetUserLocale(w, r, "fr-FR")
	assert.NoError(t, err)

	assert.Equal(t, "fr-FR", ses.GetUserLocale(r))
	assert.Equal(t, 1, sessionStoreSpy.SaveCalls)
}

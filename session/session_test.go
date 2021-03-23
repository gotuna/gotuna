package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/test/assert"
	"github.com/alcalbg/gotdd/test/doubles"
)

func TestReadingUserSIDFromEmptyStore(t *testing.T) {

	request := &http.Request{}
	sessionStoreSpy := doubles.NewSessionStoreSpy(session.GuestSID)
	ses := session.NewSession(sessionStoreSpy)

	sid, err := ses.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, session.GuestSID)
}

func TestSaveUserSIDAndRetrieve(t *testing.T) {

	request := &http.Request{}
	response := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewSessionStoreSpy(session.GuestSID)
	ses := session.NewSession(sessionStoreSpy)

	err := ses.SetUserSID(response, request, "333")
	assert.NoError(t, err)

	sid, err := ses.GetUserSID(request)
	assert.NoError(t, err)
	assert.Equal(t, sid, "333")
	assert.Equal(t, sessionStoreSpy.SaveCalls, 1)
}

func TestDestroyActiveSession(t *testing.T) {

	request := &http.Request{}
	response := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewSessionStoreSpy("333")
	ses := session.NewSession(sessionStoreSpy)

	sid, _ := ses.GetUserSID(request)
	assert.Equal(t, sid, "333")

	ses.DestroySession(response, request)

	sid, err := ses.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, session.GuestSID)
	assert.Equal(t, sessionStoreSpy.SaveCalls, 1)

}

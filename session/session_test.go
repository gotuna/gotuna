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
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(session.GuestSID)
	ses := session.NewSession(sessionStoreSpy)

	assert.Equal(t, ses.IsGuest(request), true)

	sid, err := ses.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, session.GuestSID)

	assert.Equal(t, sessionStoreSpy.SaveCalls, 0)
}

func TestSaveUserSIDAndRetrieve(t *testing.T) {

	request := &http.Request{}
	response := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(session.GuestSID)
	ses := session.NewSession(sessionStoreSpy)

	err := ses.SetUserSID(response, request, doubles.UserStub().SID)
	assert.NoError(t, err)

	sid, err := ses.GetUserSID(request)
	assert.NoError(t, err)
	assert.Equal(t, sid, doubles.UserStub().SID)
	assert.Equal(t, ses.IsGuest(request), false)
	assert.Equal(t, sessionStoreSpy.SaveCalls, 1)
}

func TestDestroyActiveSession(t *testing.T) {

	request := &http.Request{}
	response := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)
	ses := session.NewSession(sessionStoreSpy)

	sid, _ := ses.GetUserSID(request)
	assert.Equal(t, sid, doubles.UserStub().SID)

	ses.DestroySession(response, request)

	sid, err := ses.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, session.GuestSID)
	assert.Equal(t, sessionStoreSpy.SaveCalls, 1)
	assert.Equal(t, sessionStoreSpy.Session.Options.MaxAge, -1)
}

func TestFlashMessages(t *testing.T) {

	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(doubles.UserStub().SID)
	ses := session.NewSession(sessionStoreSpy)
	request := &http.Request{}
	response := httptest.NewRecorder()

	// request1: no flash messages
	messages, err := ses.Flashes(response, request)
	assert.NoError(t, err)
	assert.Equal(t, len(messages), 0)

	// request2: add flash messages
	messages, err = ses.Flashes(response, request)
	ses.AddFlash(response, request, "flash message one")
	ses.AddFlash(response, request, "flash message two")

	// request3: pop flash messages
	messages, err = ses.Flashes(response, request)
	assert.NoError(t, err)
	assert.Equal(t, len(messages), 2)

	// request4: no flash messages
	messages, err = ses.Flashes(response, request)
	assert.NoError(t, err)
	assert.Equal(t, len(messages), 0)
}

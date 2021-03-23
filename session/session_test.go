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
	sessionStoreSpy := doubles.NewSessionStoreSpy(session.GuestSID)
	ses := session.NewSession(sessionStoreSpy)
	response := httptest.NewRecorder()

	err := ses.SetUserSID(response, request, "333")
	assert.NoError(t, err)

	sid, err := ses.GetUserSID(request)
	assert.NoError(t, err)
	assert.Equal(t, sid, "333")
}

func TestDestroyActiveSession(t *testing.T) {

	request := &http.Request{}
	ses := session.NewSession(doubles.NewSessionStoreSpy(session.GuestSID))
	response := httptest.NewRecorder()

	ses.SetUserSID(response, request, "333")
	sid, _ := ses.GetUserSID(request)
	assert.Equal(t, sid, "333")

	ses.DestroySession(request)

	sid, err := ses.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, session.GuestSID)

}

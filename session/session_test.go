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
	session := session.NewSession(doubles.NewSessionStore(""))

	sid, err := session.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, "")
}

func TestSaveUserSIDAndRetrieve(t *testing.T) {

	request := &http.Request{}
	session := session.NewSession(doubles.NewSessionStore(""))
	response := httptest.ResponseRecorder{}

	err := session.SetUserSID(&response, request, "333")
	assert.NoError(t, err)

	sid, err := session.GetUserSID(request)
	assert.NoError(t, err)
	assert.Equal(t, sid, "333")
}

func TestDestroyActiveSession(t *testing.T) {

	request := &http.Request{}
	session := session.NewSession(doubles.NewSessionStore(""))
	response := httptest.ResponseRecorder{}

	session.SetUserSID(&response, request, "333")
	sid, _ := session.GetUserSID(request)
	assert.Equal(t, sid, "333")

	session.DestroySession(request)

	sid, err := session.GetUserSID(request)
	assert.Error(t, err)
	assert.Equal(t, sid, "")

}

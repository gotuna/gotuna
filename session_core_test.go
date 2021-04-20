package gotuna_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/test/assert"
	"github.com/gotuna/gotuna/test/doubles"
)

func TestStoreAndRetrieveData(t *testing.T) {

	t.Run("test storing, retrieving, and deleting a simple string", func(t *testing.T) {
		r := &http.Request{}
		w := httptest.NewRecorder()
		sessionStoreSpy := doubles.NewGorillaSessionStoreSpy("")
		ses := gotuna.NewSession(sessionStoreSpy, "test")

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
		sessionStoreSpy := doubles.NewGorillaSessionStoreSpy("")
		ses := gotuna.NewSession(sessionStoreSpy, "test")

		value, err := ses.Get(r, "test")
		assert.Equal(t, gotuna.ErrNoValueForThisKey, err)
		assert.Equal(t, "", value)
	})

}

func TestDestroyActiveSession(t *testing.T) {

	testUser := doubles.MemUser1

	r := &http.Request{}
	w := httptest.NewRecorder()
	sessionStoreSpy := doubles.NewGorillaSessionStoreSpy(testUser.GetID())
	ses := gotuna.NewSession(sessionStoreSpy, "test")

	id, err := ses.GetUserID(r)
	assert.NoError(t, err)
	assert.Equal(t, testUser.GetID(), id)

	ses.Destroy(w, r)

	id, err = ses.GetUserID(r)
	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, 1, sessionStoreSpy.SaveCalls)
	assert.Equal(t, -1, sessionStoreSpy.Session.Options.MaxAge)
}

func TestTryToUseInvalidSession(t *testing.T) {

	r := &http.Request{}
	w := httptest.NewRecorder()
	store := sessions.NewCookieStore([]byte("some key"))
	ses := gotuna.NewSession(store, "bad(((***")

	err := ses.Destroy(w, r)
	assert.Error(t, err)

	err = ses.Delete(w, r, "some key")
	assert.Error(t, err)

	_, err = ses.Get(r, "some key")
	assert.Error(t, err)

	err = ses.Put(w, r, "some key", "some value")
	assert.Error(t, err)
}

func TestSessionWillPanicOnBadSessionStore(t *testing.T) {
	defer func() {
		recover()
	}()

	gotuna.NewSession(nil, "test")

	t.Errorf("templating engine should panic")
}

func TestSessionWillPanicOnBadSessionName(t *testing.T) {
	defer func() {
		recover()
	}()

	gotuna.NewSession(doubles.NewGorillaSessionStoreSpy(""), "")

	t.Errorf("templating engine should panic")
}

func TestTypeToStringAndBackToType(t *testing.T) {

	type testType struct {
		Str string
		Bl  bool
	}

	val := testType{
		Str: "test string",
		Bl:  true,
	}

	s, err := gotuna.TypeToString(val)
	assert.NoError(t, err)
	assert.Equal(t, `{"Str":"test string","Bl":true}`, s)

	r := testType{}
	err = gotuna.TypeFromString(s, &r)
	assert.NoError(t, err)
	assert.Equal(t, "test string", r.Str)
	assert.Equal(t, true, r.Bl)
}

func TestTypeToStringWithUnsupportedType(t *testing.T) {

	type testType struct {
		Str string
		Fn  func(string)
	}

	val := testType{
		Str: "test string",
		Fn:  func(a string) {},
	}

	_, err := gotuna.TypeToString(val)
	assert.Error(t, err)
}

func TestTypeFromGarbageString(t *testing.T) {
	err := gotuna.TypeFromString("garbage===", struct{}{})
	assert.Error(t, err)
}

package stubs

import (
	"io"
	"log"
	"net/http"

	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/sessions"
)

func NewLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func NewSessionStore(r *http.Request, userSID string) sessions.Store {

	userSession := sessions.NewSession(nil, "stub_session")
	userSession.Values[session.UserSIDKey] = userSID

	return &stubSessionStore{userSession}
}

type stubSessionStore struct {
	userSession *sessions.Session
}

func (session *stubSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return session.userSession, nil
}

func (session *stubSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return session.userSession, nil
}

func (session *stubSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

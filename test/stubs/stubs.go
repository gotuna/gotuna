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

	userSession := sessions.NewSession(nil, "")
	userSession.Values[session.UserSIDKey] = userSID

	return &sessionStore{userSession}
}

type sessionStore struct {
	userSession *sessions.Session
}

func (stub *sessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return stub.userSession, nil
}

func (stub *sessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return stub.userSession, nil
}

func (stub *sessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

package doubles

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

func NewSessionStoreSpy(r *http.Request, userSID string) *SessionStoreSpy {

	userSession := sessions.NewSession(nil, "")
	userSession.Values[session.UserSIDKey] = userSID

	return &SessionStoreSpy{
		userSession: userSession,
		GetCalls:    0,
	}
}

type SessionStoreSpy struct {
	userSession *sessions.Session
	GetCalls    int
}

func (stub *SessionStoreSpy) Get(r *http.Request, name string) (*sessions.Session, error) {
	stub.GetCalls++
	return stub.userSession, nil
}

func (stub *SessionStoreSpy) New(r *http.Request, name string) (*sessions.Session, error) {
	return stub.userSession, nil
}

func (stub *SessionStoreSpy) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

package doubles

import (
	"io"
	"log"
	"net/http"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/sessions"
)

func NewLoggerStub() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func NewSessionStoreSpy(userSID string) *SessionStoreSpy {
	userSession := sessions.NewSession(&SessionStoreSpy{}, "")
	userSession.Values[session.UserSIDKey] = userSID

	return &SessionStoreSpy{Session: userSession}
}

type SessionStoreSpy struct {
	Session   *sessions.Session
	GetCalls  int
	NewCalls  int
	SaveCalls int
}

func (stub *SessionStoreSpy) Get(r *http.Request, name string) (*sessions.Session, error) {
	stub.GetCalls++
	return stub.Session, nil
}

func (stub *SessionStoreSpy) New(r *http.Request, name string) (*sessions.Session, error) {
	stub.NewCalls++
	return stub.Session, nil
}

func (stub *SessionStoreSpy) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	stub.SaveCalls++
	return nil
}

func NewUserRepositoryStub(user app.User) UserRepositoryStub {
	return UserRepositoryStub{user}
}

type UserRepositoryStub struct {
	user app.User
}

func (u UserRepositoryStub) GetUserByEmail(email string) (app.User, error) {
	return u.user, nil
}

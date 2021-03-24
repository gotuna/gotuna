package doubles

import (
	"errors"
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

func NewGorillaSessionStoreSpy(userSID string) *GorillaSessionStoreSpy {
	userSession := sessions.NewSession(&GorillaSessionStoreSpy{}, "")
	userSession.Values[session.UserSIDKey] = userSID

	return &GorillaSessionStoreSpy{session: userSession}
}

type GorillaSessionStoreSpy struct {
	session   *sessions.Session
	GetCalls  int
	NewCalls  int
	SaveCalls int
}

func (spy *GorillaSessionStoreSpy) Get(r *http.Request, name string) (*sessions.Session, error) {
	spy.GetCalls++
	return spy.session, nil
}

func (spy *GorillaSessionStoreSpy) New(r *http.Request, name string) (*sessions.Session, error) {
	spy.NewCalls++
	return spy.session, nil
}

func (spy *GorillaSessionStoreSpy) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	spy.SaveCalls++
	return nil
}

func NewUserRepositoryStub(user app.User) UserRepositoryStub {
	return UserRepositoryStub{user}
}

type UserRepositoryStub struct {
	user app.User
}

func (u UserRepositoryStub) GetUserByEmail(email string) (app.User, error) {
	if u.user.Email != email {
		return app.User{}, errors.New("no user")
	}
	return u.user, nil
}

func NewServerStub() *app.Server {
	return app.NewServer(
		NewLoggerStub(),
		session.NewSession(NewGorillaSessionStoreSpy(session.GuestSID)),
		NewUserRepositoryStub(UserStub()),
	)
}

func NewServerWithCookieStoreStub() *app.Server {
	return app.NewServer(
		NewLoggerStub(),
		session.NewSession(sessions.NewCookieStore([]byte("abc"))),
		NewUserRepositoryStub(UserStub()),
	)
}

func UserStub() app.User {
	return app.User{
		SID:          "123",
		Email:        "john@example.com",
		PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm", // pass123
	}
}
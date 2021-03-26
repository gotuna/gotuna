package doubles

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/templating"
	"github.com/gorilla/sessions"
)

func NewLoggerStub() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func NewGorillaSessionStoreSpy(userSID string) *GorillaSessionStoreSpy {
	userSession := sessions.NewSession(&GorillaSessionStoreSpy{}, "")
	userSession.Values[session.UserSIDKey] = userSID

	return &GorillaSessionStoreSpy{Session: userSession}
}

type GorillaSessionStoreSpy struct {
	Session   *sessions.Session
	GetCalls  int
	NewCalls  int
	SaveCalls int
}

func (spy *GorillaSessionStoreSpy) Get(r *http.Request, name string) (*sessions.Session, error) {
	spy.GetCalls++
	return spy.Session, nil
}

func (spy *GorillaSessionStoreSpy) New(r *http.Request, name string) (*sessions.Session, error) {
	spy.NewCalls++
	return spy.Session, nil
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
		NewFileSystemStub(nil),
		session.NewSession(NewGorillaSessionStoreSpy(session.GuestSID)),
		NewUserRepositoryStub(UserStub()),
	)
}

func NewServerWithCookieStoreStub() *app.Server {
	return app.NewServer(
		NewLoggerStub(),
		NewFileSystemStub(nil),
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

func NewFileSystemStub(files map[string][]byte) *filesystemStub {
	return &filesystemStub{
		files: files,
	}
}

type filesystemStub struct {
	files map[string][]byte
}

func (f *filesystemStub) Open(name string) (fs.File, error) {
	tmpfile, err := ioutil.TempFile("", "fsdemo")
	if err != nil {
		log.Fatal(err)
	}

	contents, ok := f.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	tmpfile.Write([]byte(contents))
	tmpfile.Seek(0, 0)

	return tmpfile, nil
}

var StubTemplate = `{{define "app"}}{{end}}`

func NewStubTemplatingEngine(template string, session *session.Session) templating.TemplatingEngine {
	return templating.GetEngine(i18n.NewTranslator(i18n.En), session).
		MountFS(
			NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(template),
				}))
}

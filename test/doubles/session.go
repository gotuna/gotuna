package doubles

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gotuna/gotuna"
)

func NewGorillaSessionStoreSpy(userID string) *storeSpy {
	userSession := sessions.NewSession(&storeSpy{}, "")
	userSession.Values[gotuna.UserIDKey] = userID

	return &storeSpy{Session: userSession}
}

// implements gorilla.Store interface
type storeSpy struct {
	Session   *sessions.Session
	GetCalls  int
	NewCalls  int
	SaveCalls int
}

func (spy *storeSpy) Get(r *http.Request, name string) (*sessions.Session, error) {
	spy.GetCalls++
	return spy.Session, nil
}

func (spy *storeSpy) New(r *http.Request, name string) (*sessions.Session, error) {
	spy.NewCalls++
	return spy.Session, nil
}

func (spy *storeSpy) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	spy.SaveCalls++
	return nil
}

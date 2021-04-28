package doubles

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gotuna/gotuna"
)

// NewGorillaSessionStoreSpy return a new gorilla.Store spy
func NewGorillaSessionStoreSpy(userID string) *StoreSpy {
	userSession := sessions.NewSession(&StoreSpy{}, "")
	userSession.Values[gotuna.UserIDKey] = userID

	return &StoreSpy{Session: userSession}
}

// StoreSpy implements gorilla.Store interface
type StoreSpy struct {
	Session   *sessions.Session
	GetCalls  int
	NewCalls  int
	SaveCalls int
}

// Get counts the Get calls
func (spy *StoreSpy) Get(r *http.Request, name string) (*sessions.Session, error) {
	spy.GetCalls++
	return spy.Session, nil
}

// New counts the New calls
func (spy *StoreSpy) New(r *http.Request, name string) (*sessions.Session, error) {
	spy.NewCalls++
	return spy.Session, nil
}

// Save counts the Save calls
func (spy *StoreSpy) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	spy.SaveCalls++
	return nil
}

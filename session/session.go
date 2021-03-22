package session

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const UserSIDKey = "user_sid"
const SessionName = "app_session"

type Session struct {
	store sessions.Store
}

// NewSession returns new session with store
func NewSession(store sessions.Store) *Session {
	if store == nil {
		panic("Must supply a session store")
	}

	return &Session{store: store}
}

func (s Session) SetUserSID(w http.ResponseWriter, r *http.Request, sid string) error {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return errors.New("Cannot get session from the store")
	}

	session.Values[UserSIDKey] = sid

	if err = session.Save(r, w); err != nil {
		return errors.New("Cannot store to session")
	}

	return nil
}

func (s Session) GetUserSID(r *http.Request) (string, error) {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return "", errors.New("Cannot get a session from the store")
	}

	sid, ok := session.Values[UserSIDKey].(string)
	if !ok || sid == "" {
		return "", errors.New("No user in the session")
	}

	return sid, nil
}

func (s Session) DestroySession(r *http.Request) error {
	session, err := s.store.Get(r, SessionName)
	if err != nil {
		return errors.New("Cannot get a session from the store")
	}

	delete(session.Values, UserSIDKey)

	return nil
}

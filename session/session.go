package session

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const GuestSID = ""
const UserSIDKey = "user_sid"
const sessionName = "app_session"

type Session struct {
	Store sessions.Store
}

// NewSession returns new session with requested store
func NewSession(store sessions.Store) *Session {
	if store == nil {
		panic("Must supply a valid session store")
	}

	return &Session{Store: store}
}

func (s Session) SetUserSID(w http.ResponseWriter, r *http.Request, sid string) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("Cannot get session from the store")
	}

	session.Values[UserSIDKey] = sid

	if err = s.Store.Save(r, w, session); err != nil {
		return errors.New("Cannot store to session")
	}

	return nil
}

func (s Session) GetUserSID(r *http.Request) (string, error) {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return GuestSID, errors.New("Cannot get session from the store")
	}

	sid, ok := session.Values[UserSIDKey].(string)
	if !ok || sid == GuestSID {
		return GuestSID, errors.New("No user in the session")
	}

	return sid, nil
}

func (s Session) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, sessionName)
	if err != nil {
		return errors.New("Cannot get session from the store")
	}

	delete(session.Values, UserSIDKey)
	session.Options.MaxAge = -1

	if err = s.Store.Save(r, w, session); err != nil {
		return errors.New("Cannot store to session")
	}

	return nil
}

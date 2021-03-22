package session

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const UserSIDKey = "user_sid"
const SessionName = "app_session"

func SetUserSID(w http.ResponseWriter, r *http.Request, store sessions.Store, sid string) error {
	session, err := store.Get(r, SessionName)
	if err != nil {
		return errors.New("Cannot get session from the store")
	}

	session.Values[UserSIDKey] = sid

	if err = session.Save(r, w); err != nil {
		return errors.New("Cannot store to session")
	}

	return nil
}

func GetUserSID(r *http.Request, store sessions.Store) (string, error) {
	session, err := store.Get(r, SessionName)
	if err != nil {
		return "", errors.New("Cannot get a session from the store")
	}

	sid, ok := session.Values[UserSIDKey].(string)
	if !ok || sid == "" {
		return "", errors.New("No user in the session")
	}

	return sid, nil
}

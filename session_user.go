package gotdd

import (
	"errors"
	"net/http"
)

const UserIDKey = "_user_id"
const UserLocaleKey = "_user_locale"

func (s Session) IsGuest(r *http.Request) bool {
	id, err := s.GetUserID(r)
	if err != nil || id == "" {
		return true
	}
	return false
}

func (s Session) SetUserID(w http.ResponseWriter, r *http.Request, id string) error {
	return s.Put(w, r, UserIDKey, id)
}

func (s Session) GetUserID(r *http.Request) (string, error) {
	id, err := s.Get(r, UserIDKey)
	if err != nil || id == "" {
		return "", errors.New("no user in the session")
	}

	return id, nil
}

func (s Session) SetUserLocale(w http.ResponseWriter, r *http.Request, id string) error {
	return s.Put(w, r, UserLocaleKey, id)
}

func (s Session) GetUserLocale(r *http.Request) string {
	locale, err := s.Get(r, UserLocaleKey)
	if err != nil {
		return "" // TODO: default locale
	}

	return locale
}

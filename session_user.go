package gotdd

import (
	"errors"
	"net/http"
)

const GuestSID = ""
const UserSIDKey = "_user_sid"
const UserLocaleKey = "_user_locale"

func (s Session) IsGuest(r *http.Request) bool {
	sid, err := s.GetUserSID(r)
	if err != nil || sid == GuestSID {
		return true
	}

	return false
}

func (s Session) SetUserSID(w http.ResponseWriter, r *http.Request, sid string) error {
	return s.Put(w, r, UserSIDKey, sid)
}

func (s Session) GetUserSID(r *http.Request) (string, error) {
	sid, err := s.Get(r, UserSIDKey)
	if err != nil || sid == GuestSID {
		return GuestSID, errors.New("no user in the session")
	}

	return sid, nil
}

func (s Session) SetUserLocale(w http.ResponseWriter, r *http.Request, sid string) error {
	return s.Put(w, r, UserLocaleKey, sid)
}

func (s Session) GetUserLocale(r *http.Request) string {
	locale, err := s.Get(r, UserLocaleKey)
	if err != nil {
		return "" // TODO: default locale
	}

	return locale
}

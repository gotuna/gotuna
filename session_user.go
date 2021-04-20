package gotuna

import (
	"errors"
	"net/http"
)

const (
	// UserIDKey is used as session key to store the current's user unique ID.
	UserIDKey = "_user_id"
	// UserLocaleKey is used as session key for the current's user locale settings.
	UserLocaleKey = "_user_locale"
)

// IsGuest checks if current user is not logged in into the app.
func (s Session) IsGuest(r *http.Request) bool {
	id, err := s.GetUserID(r)
	if err != nil || id == "" {
		return true
	}
	return false
}

// SetUserID will save the current user's unique ID into the session.
func (s Session) SetUserID(w http.ResponseWriter, r *http.Request, id string) error {
	return s.Put(w, r, UserIDKey, id)
}

// GetUserID retrieves the current user's unique ID from the session.
func (s Session) GetUserID(r *http.Request) (string, error) {
	id, err := s.Get(r, UserIDKey)
	if err != nil || id == "" {
		return "", errors.New("no user in the session")
	}

	return id, nil
}

// SetUserLocale will store the user's locale string into the session.
func (s Session) SetUserLocale(w http.ResponseWriter, r *http.Request, id string) error {
	return s.Put(w, r, UserLocaleKey, id)
}

// GetUserLocale retrieves the current user's locale string from the session.
func (s Session) GetUserLocale(r *http.Request) string {
	locale, err := s.Get(r, UserLocaleKey)
	if err != nil {
		return "" // TODO: default locale
	}

	return locale
}

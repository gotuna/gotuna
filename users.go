package gotuna

import (
	"net/http"
)

// The User interface for providing the standard user entity.
type User interface {
	GetID() string
}

// The UserRepository will provide a way to retrieve and authenticate users.
//
// The user repository is database-agnostic. You can keep your users in MySQL,
// MongoDB, LDAP or in a simple JSON file. By implementing this interface, you
// are providing the way for the app to authenticate, and retrieve the specific
// user entity by using the unique user ID.
type UserRepository interface {
	GetUserByID(id string) (User, error)
	Authenticate(w http.ResponseWriter, r *http.Request) (User, error)
}

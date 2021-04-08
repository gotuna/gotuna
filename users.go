package gotuna

import (
	"net/http"
)

type User interface {
	GetID() string
}

type UserRepository interface {
	Authenticate(w http.ResponseWriter, r *http.Request) (User, error)
	GetUserByID(id string) (User, error)
}

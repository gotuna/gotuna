package gotdd

import "net/http"

type User interface {
	GetID() string
}

type UserRepository interface {
	Authenticate(w http.ResponseWriter, r *http.Request) (User, error)
	GetByID(id string) (User, error)
}

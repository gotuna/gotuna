package gotuna

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// InMemoryUser is a sample User entity implementation with some
// sample fields provided like Name and Email.
//
// Password is stored as plain-text for simplicity. In real life, you should
// probably use crypto/bcrypt library and store hashed representation and
// compare passwords in Authenticate method with bcrypt.CompareHashAndPassword
type InMemoryUser struct {
	ID       string
	Email    string
	Name     string
	Password string
}

// GetID will return a unique ID for every specific user.
func (u InMemoryUser) GetID() string {
	return u.ID
}

// NewInMemoryUserRepository returns a sample in-memory UserRepository
// implementation which can be used in tests or sample apps.
func NewInMemoryUserRepository(users []InMemoryUser) UserRepository {
	return inMemoryUserRepository{users}
}

type inMemoryUserRepository struct {
	users []InMemoryUser
}

func (u inMemoryUserRepository) Authenticate(w http.ResponseWriter, r *http.Request) (User, error) {

	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	Password := r.FormValue("password")

	if email == "" {
		return InMemoryUser{}, errors.New("this field is required")
	}
	if Password == "" {
		return InMemoryUser{}, errors.New("this field is required")
	}

	found, err := u.getUserByEmail(email)
	if err != nil {
		return InMemoryUser{}, fmt.Errorf("cannot find user with this email %v", err)
	}

	// in real life this should be bcrypt.CompareHashAndPassword
	if found.Password != Password {
		return InMemoryUser{}, fmt.Errorf("passwords don't match %v", err)
	}

	return found, nil
}

func (u inMemoryUserRepository) GetUserByID(id string) (User, error) {
	for _, user := range u.users {
		if user.ID == id {
			return user, nil
		}
	}

	return InMemoryUser{}, errors.New("user not found")
}

func (u inMemoryUserRepository) getUserByEmail(email string) (InMemoryUser, error) {
	for _, user := range u.users {
		if user.Email == email {
			return user, nil
		}
	}

	return InMemoryUser{}, errors.New("user not found")
}

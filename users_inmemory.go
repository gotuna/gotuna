package gotuna

import (
	"errors"
	"net/http"
	"strings"
)

var (
	// ErrWrongPassword is thrown on bad password
	ErrWrongPassword = constError("wrong password")
	// ErrCannotFindUser is thrown when we cannot match user
	ErrCannotFindUser = constError("user not found")
	// ErrRequiredField is thrown when requestfield is required
	ErrRequiredField = constError("this field is required")
	// ErrNotUnique is thrown when we try to add new user with existing ID
	ErrNotUnique = constError("user with this ID already exists")
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
	return &InMemoryUserRepository{users}
}

// InMemoryUserRepository is a sample in-memory UserRepository
type InMemoryUserRepository struct {
	Users []InMemoryUser
}

// Authenticate tries to authenticate the user when submitting the login form
func (u InMemoryUserRepository) Authenticate(w http.ResponseWriter, r *http.Request) (User, error) {

	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	Password := r.FormValue("password")

	if email == "" {
		return InMemoryUser{}, ErrRequiredField
	}
	if Password == "" {
		return InMemoryUser{}, ErrRequiredField
	}

	found, err := u.getUserByEmail(email)
	if err != nil {
		return InMemoryUser{}, ErrCannotFindUser
	}

	// in real life this should be bcrypt.CompareHashAndPassword
	if found.Password != Password {
		return InMemoryUser{}, ErrWrongPassword
	}

	return found, nil
}

// GetUserByID will try to find and return the user with this ID from the repository
func (u InMemoryUserRepository) GetUserByID(id string) (User, error) {
	for _, user := range u.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return InMemoryUser{}, ErrCannotFindUser
}

func (u InMemoryUserRepository) getUserByEmail(email string) (InMemoryUser, error) {
	for _, user := range u.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return InMemoryUser{}, errors.New("user not found")
}

// AddUser will add user to the repository or return error if ID is not unique or blank
func (u *InMemoryUserRepository) AddUser(user InMemoryUser) error {
	if _, err := u.GetUserByID(user.ID); user.ID == "" || err != ErrCannotFindUser {
		return ErrNotUnique
	}
	u.Users = append(u.Users, user)
	return nil
}

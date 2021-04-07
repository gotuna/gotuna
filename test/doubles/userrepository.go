package doubles

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd"
)

var FakeUser1 = fakeUserStub{
	databaseID: "123",
	Email:      "john@example.com",
	Name:       "John",
	password:   "pass123",
}

var FakeUser2 = fakeUserStub{
	databaseID: "456",
	Email:      "bob@example.com",
	Name:       "Bob",
	password:   "bobby5",
}

type fakeUserStub struct {
	databaseID string
	Email      string
	Name       string
	password   string
}

func (u fakeUserStub) GetID() string {
	return u.databaseID
}

func NewUserRepositoryStub() *userRepositoryStub {
	return &userRepositoryStub{
		users: []fakeUserStub{
			FakeUser1,
			FakeUser2,
		},
	}
}

type userRepositoryStub struct {
	users []fakeUserStub
}

func (u userRepositoryStub) Authenticate(w http.ResponseWriter, r *http.Request) (gotdd.User, error) {

	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	if email == "" {
		return fakeUserStub{}, errors.New("this field is required")
	}
	if password == "" {
		return fakeUserStub{}, errors.New("this field is required")
	}

	found, err := u.getUserByEmail(email)
	if err != nil {
		return fakeUserStub{}, fmt.Errorf("cannot find user with this email %v", err)
	}

	// in real life this should be bcrypt.CompareHashAndPassword
	if found.password != password {
		return fakeUserStub{}, fmt.Errorf("passwords don't match %v", err)
	}

	return found, nil
}

func (u userRepositoryStub) GetUserByID(id string) (gotdd.User, error) {
	for _, user := range u.users {
		if user.databaseID == id {
			return user, nil
		}
	}

	return fakeUserStub{}, errors.New("user not found")
}

func (u userRepositoryStub) getUserByEmail(email string) (fakeUserStub, error) {
	for _, user := range u.users {
		if user.Email == email {
			return user, nil
		}
	}

	return fakeUserStub{}, errors.New("user not found")
}

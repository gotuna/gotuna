package doubles

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd"
)

var FakeUser1 = FakeUserStub{
	databaseID: "123",
	Email:      "john@example.com",
	password:   "pass123",
}

var FakeUser2 = FakeUserStub{
	databaseID: "456",
	Email:      "bob@example.com",
	password:   "bobby5",
}

type FakeUserStub struct {
	databaseID string
	Email      string
	password   string
}

func (u FakeUserStub) GetID() string {
	return u.databaseID
}

func NewUserRepositoryStub() *userRepositoryStub {
	return &userRepositoryStub{
		users: []FakeUserStub{
			FakeUser1,
			FakeUser2,
		},
	}
}

type userRepositoryStub struct {
	users []FakeUserStub
}

func (u userRepositoryStub) Authenticate(w http.ResponseWriter, r *http.Request) (gotdd.User, error) {

	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	if email == "" {
		return FakeUserStub{}, errors.New("this field is required")
	}
	if password == "" {
		return FakeUserStub{}, errors.New("this field is required")
	}

	found, err := u.getUserByEmail(email)
	if err != nil {
		return FakeUserStub{}, fmt.Errorf("cannot find user with this email %v", err)
	}

	// in real life this should be bcrypt.CompareHashAndPassword
	if found.password != password {
		return FakeUserStub{}, fmt.Errorf("passwords don't match %v", err)
	}

	return found, nil
}

func (u userRepositoryStub) GetUserByID(id string) (gotdd.User, error) {
	for _, user := range u.users {
		if user.databaseID == id {
			return user, nil
		}
	}

	return FakeUserStub{}, errors.New("user not found")
}

func (u userRepositoryStub) getUserByEmail(email string) (FakeUserStub, error) {
	for _, user := range u.users {
		if user.Email == email {
			return user, nil
		}
	}

	return FakeUserStub{}, errors.New("user not found")
}

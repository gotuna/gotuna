package doubles

import (
	"errors"
	"fmt"

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
	users         []FakeUserStub
	inputEmail    string
	inputPassword string
}

func (u userRepositoryStub) Authenticate() (gotdd.User, error) {
	found, err := u.getUserByEmail(u.inputEmail)
	if err != nil {
		return FakeUserStub{}, err
	}

	// in real life this should be bcrypt.CompareHashAndPassword
	if found.password != u.inputPassword {
		return FakeUserStub{}, fmt.Errorf("passwords don't match %v", err)
	}

	return found, nil
}

func (u *userRepositoryStub) Set(key string, value interface{}) gotdd.UserRepository {
	if key == "email" {
		u.inputEmail = value.(string)
	}
	if key == "password" {
		u.inputPassword = value.(string)
	}
	return u
}

func (u userRepositoryStub) GetByID(id string) (gotdd.User, error) {
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

package doubles

import (
	"errors"
	"fmt"

	"github.com/alcalbg/gotdd"
)

func UserStub() gotdd.User {
	return gotdd.User{
		SID:          "123",
		Email:        "john@example.com",
		PasswordHash: "pass123",
		Locale:       "en-US",
	}
}

func NewUserRepositoryStub() *userRepositoryStub {
	return &userRepositoryStub{
		users: []gotdd.User{UserStub()},
	}
}

type userRepositoryStub struct {
	users         []gotdd.User
	inputEmail    string
	inputPassword string
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

func (u userRepositoryStub) Authenticate() (gotdd.User, error) {
	user, err := u.getUserByEmail()
	if err != nil {
		return gotdd.User{}, err
	}

	// this should be bcrypt.CompareHashAndPassword in real life
	if user.PasswordHash != u.inputPassword {
		return gotdd.User{}, fmt.Errorf("passwords don't match %v", err)
	}

	return user, nil
}

func (u userRepositoryStub) getUserByEmail() (gotdd.User, error) {
	for _, user := range u.users {
		if user.Email == u.inputEmail {
			return user, nil
		}
	}

	return gotdd.User{}, errors.New("user not found")
}

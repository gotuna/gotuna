package doubles

import (
	"errors"
	"fmt"

	"github.com/alcalbg/gotdd/models"
)

func UserStub() models.User {
	return models.User{
		SID:          "123",
		Email:        "john@example.com",
		PasswordHash: "pass123",
	}
}

func NewUserRepositoryStub() *userRepositoryStub {
	return &userRepositoryStub{
		users: []models.User{UserStub()},
	}
}

type userRepositoryStub struct {
	users         []models.User
	inputEmail    string
	inputPassword string
}

func (u *userRepositoryStub) Set(key string, value interface{}) models.UserRepository {
	if key == "email" {
		u.inputEmail = value.(string)
	}
	if key == "password" {
		u.inputPassword = value.(string)
	}
	return u
}

func (u userRepositoryStub) Authenticate() (models.User, error) {
	user, err := u.getUserByEmail()
	if err != nil {
		return models.User{}, err
	}

	// this should be bcrypt.CompareHashAndPassword in real life
	if user.PasswordHash != u.inputPassword {
		return models.User{}, fmt.Errorf("passwords don't match %v", err)
	}

	return user, nil
}

func (u userRepositoryStub) getUserByEmail() (models.User, error) {
	for _, user := range u.users {
		if user.Email == u.inputEmail {
			return user, nil
		}
	}

	return models.User{}, errors.New("user not found")
}

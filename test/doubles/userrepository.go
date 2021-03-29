package doubles

import (
	"errors"

	"github.com/alcalbg/gotdd/models"
)

func UserStub() models.User {
	return models.User{
		SID:          "123",
		Email:        "john@example.com",
		PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm", // pass123
	}
}

func NewUserRepositoryStub(user models.User) userRepositoryStub {
	return userRepositoryStub{user}
}

// implements UserRepository interface
type userRepositoryStub struct {
	user models.User
}

func (u userRepositoryStub) GetUserByEmail(email string) (models.User, error) {
	if u.user.Email != email {
		return models.User{}, errors.New("no user")
	}
	return u.user, nil
}

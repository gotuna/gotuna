package models

import "errors"

type User struct {
	SID          string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	GetUserByEmail(email string) (User, error)
}

func NewInMemoryUserRepository() UserRepository {
	repo := &inMemoryUserRepository{}
	repo.users = []User{
		User{SID: "1", Email: "alcalbg@gmail.com", PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm"}, // pass123
		User{SID: "2", Email: "admin@example.com", PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm"}, // pass123
	}
	return repo
}

type inMemoryUserRepository struct {
	users []User
}

func (repo inMemoryUserRepository) GetUserByEmail(email string) (User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")

}

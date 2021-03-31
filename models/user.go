package models

type User struct {
	SID          string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	Authenticate() (User, error)
	Set(key string, value interface{}) UserRepository
}

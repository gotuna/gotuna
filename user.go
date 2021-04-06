package gotdd

type User interface {
	GetID() string
}

type UserRepository interface {
	Authenticate() (User, error)
	Set(key string, value interface{}) UserRepository
	GetByID(id string) (User, error)
}

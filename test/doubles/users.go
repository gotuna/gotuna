package doubles

import "github.com/gotuna/gotuna"

var MemUser1 = gotuna.InMemoryUser{
	UniqueID: "123",
	Email:    "john@example.com",
	Name:     "John",
	Password: "pass123",
}

var MemUser2 = gotuna.InMemoryUser{
	UniqueID: "456",
	Email:    "bob@example.com",
	Name:     "Bob",
	Password: "bobby5",
}

func NewUserRepositoryStub() gotuna.UserRepository {
	return gotuna.NewInMemoryUserRepository([]gotuna.InMemoryUser{
		MemUser1,
		MemUser2,
	})
}

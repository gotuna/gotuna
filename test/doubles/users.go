package doubles

import (
	"github.com/alcalbg/gotdd"
)

var MemUser1 = gotdd.InMemoryUser{
	UniqueID: "123",
	Email:    "john@example.com",
	Name:     "John",
	Password: "pass123",
}

var MemUser2 = gotdd.InMemoryUser{
	UniqueID: "456",
	Email:    "bob@example.com",
	Name:     "Bob",
	Password: "bobby5",
}

func NewUserRepositoryStub() gotdd.UserRepository {
	return gotdd.NewInMemoryUserRepository([]gotdd.InMemoryUser{
		MemUser1,
		MemUser2,
	})
}

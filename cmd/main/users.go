package main

import "github.com/alcalbg/gotdd"

var User1 = gotdd.InMemoryUser{
	UniqueID: "123",
	Email:    "john@example.com",
	Name:     "John",
	Password: "pass123",
}

var User2 = gotdd.InMemoryUser{
	UniqueID: "456",
	Email:    "bob@example.com",
	Name:     "Bob",
	Password: "bobby5",
}

func NewUserRepository() gotdd.UserRepository {
	return gotdd.NewInMemoryUserRepository([]gotdd.InMemoryUser{
		User1,
		User2,
	})
}

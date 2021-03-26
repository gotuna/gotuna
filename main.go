package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/static"
	"github.com/gorilla/sessions"
)

type MemoryUserRepository struct {
	users []app.User
}

func NewMemoryUserRepository() app.UserRepository {
	repo := &MemoryUserRepository{}
	repo.users = []app.User{
		app.User{SID: "1", Email: "alcalbg@gmail.com", PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm"}, // pass123
		app.User{SID: "2", Email: "admin@example.com", PasswordHash: "$2a$10$19ogjdlTWc0dHBeC5i1qOeNP6oqwIgphXmtrpjFBt3b4ru5B5Cxfm"}, // pass123
	}
	return repo
}

func (repo *MemoryUserRepository) GetUserByEmail(email string) (app.User, error) {
	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}

	return app.User{}, errors.New("user not found")

}

func main() {

	port := ":8888"

	gorillaSessionStore := sessions.NewCookieStore([]byte(os.Getenv("APP_KEY")))
	fs := static.EmbededStatic

	srv := app.NewServer(
		log.New(os.Stdout, "", 0),
		fs,
		session.NewSession(gorillaSessionStore),
		NewMemoryUserRepository(),
	)

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, srv.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

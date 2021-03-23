package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/sessions"
)

type MemoryUserRepository struct {
	users []app.User
}

func NewMemoryUserRepository() app.UserRepository {
	repo := &MemoryUserRepository{}
	repo.users = []app.User{
		app.User{SID: "1", Email: "alcalbg@gmail.com", PasswordHash: ""},
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

	cookieStore := sessions.NewCookieStore([]byte(os.Getenv("APP_KEY")))
	srv := app.NewServer(
		log.New(os.Stdout, "", 0),
		sessions.NewSession(cookieStore, session.SessionName),
		NewMemoryUserRepository(),
	)

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, srv.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

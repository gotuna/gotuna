package server

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Logger *log.Logger
}

const get = http.MethodGet
const post = http.MethodPost
const DefaultError = "Whoops, something went wrong"

func NewServer(logger *log.Logger) *Server {
	s := &Server{}
	s.Logger = logger

	s.Router = mux.NewRouter()
	s.Router.NotFoundHandler = http.HandlerFunc(notFound)

	s.Router.Handle("/", home()).Methods(get)
	s.Router.Handle("/login", login()).Methods(get)

	//bad := func() http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		var x map[string]int
	//		x["y"] = 1 // will produce nil map panic
	//	})
	//}
	//s.Router.Handle("/bad", bad())

	s.Router.Use(s.recoverer)

	return s
}

func home() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (s Server) recoverer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				// log error and stack trace to console
				s.Logger.Printf("PANIC RECOVERED: %v", err)
				s.Logger.Println(string(debug.Stack()))

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(DefaultError))
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

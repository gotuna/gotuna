package gotdd

import "net/http"

type Server struct{}

func (srv Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {}

func NewServer() http.Handler {
	return Server{}
}

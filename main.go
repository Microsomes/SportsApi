package main

import (
	"io"
	"net/http"

	"github.com/sportsapiv1/handlers"
)

type Server struct {
	*http.Server
}

func NewServer(sm *http.ServeMux) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:    ":5001",
			Handler: sm,
		},
	}
	return s
}

func (s *Server) block() {
	select {}
}

func AskInput(b io.Writer, what []byte) {
	b.Write(what)
}

func main() {

	sm := http.NewServeMux()

	sm.HandleFunc("/", handlers.AllDocs().ServeHTTP)
	sm.HandleFunc("/leagues", handlers.AllLeagues().ServeHTTP)

	sm.HandleFunc("/languages", handlers.AllLanguages)

	s := NewServer(sm)

	s.ListenAndServe()

	s.block()

}

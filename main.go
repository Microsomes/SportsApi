package main

import (
	"io"
	"net/http"
	"os"

	"github.com/sportsapiv1/handlers"
)

type Server struct {
	*http.Server
}

func NewServer(sm *http.ServeMux, port string) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:    ":" + port,
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
	sm.HandleFunc("/language/", handlers.OneLanguage)

	var port string = ""
	if os.Getenv("PORT") == "" {
		port = "5001"
	} else {
		port = os.Getenv("PORT")
	}

	s := NewServer(sm, port)

	s.ListenAndServe()

	s.block()

}

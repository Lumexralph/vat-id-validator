// Package server contains implementation for the server,
// and routing patterns to the handlers.
package server

import (
	"github.com/Lumexralph/vat-id-validator/validator"
	"log"
	"net/http"
	"os"
)

const DefaultPort = "3000"

type Server struct {
	vatChecker validator.VATIDChecker
}

func NewServer(vatChecker validator.VATIDChecker) *Server {
	return &Server{
		vatChecker: vatChecker,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.indexHandler)
	mux.HandleFunc("/vatid/validate", s.vatIDHandler)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = DefaultPort
	}

	log.Printf("Starting server on port:%s... \n", port)
	return http.ListenAndServe(":"+port, mux)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("German VATID Validator Microservice\n"))
}

func (s *Server) vatIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("German VATID Validator Microservice - Validate"))
}

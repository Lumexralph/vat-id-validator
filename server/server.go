// Package server contains implementation for the server,
// and routing patterns to the handlers.
package server

import (
	"encoding/json"
	"github.com/Lumexralph/vat-id-validator/validator"
	"log"
	"net/http"
	"os"
)

const DefaultPort = "3000"

type VATPost struct {
	VATNumber string `json:"vat_number"`
}

type VATPostResponse struct {
	Valid bool `json:"valid"`
}

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
	if r.Method != http.MethodPost {
		http.Error(w, "HTTP method not supported", http.StatusMethodNotAllowed)
		return
	}

	var post VATPost

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: validate empty request

	valid, err := s.vatChecker.ValidateVATID(r.Context(), post.VATNumber)
	if err != nil {
		log.Printf("error reported: %v\n", err)
		return
	}
	log.Println("valid response: ", valid)

	res := VATPostResponse{
		Valid: valid,
	}

	output, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	log.Printf("request: %q", post)
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

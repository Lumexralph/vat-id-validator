// Package server contains implementation for the server,
// and routing patterns to the handlers.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Lumexralph/vat-id-validator/validator"
)

const DefaultPort = "3000"

type vATPost struct {
	VATNumber string `json:"vat_number"`
}

type vATPostResponse struct {
	Valid bool `json:"valid"`
}

type server struct {
	vatChecker validator.VATIDChecker
}

func NewServer(vatChecker validator.VATIDChecker) *server {
	return &server{
		vatChecker: vatChecker,
	}
}

func (s *server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.indexHandler)
	mux.HandleFunc("/vatid/validate", s.vatIDHandler)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = DefaultPort
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go gracefulServerShutdown(srv, idleConnsClosed)

	log.Printf("Starting server on port:%s... \n", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return fmt.Errorf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

	return nil
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("German VATID Validator Microservice\n"))
}

func (s *server) vatIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "HTTP method not supported", http.StatusMethodNotAllowed)
		return
	}

	var post vATPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if vat := strings.ReplaceAll(post.VATNumber, " ", ""); vat == "" {
		http.Error(w, "vat_number not provided", http.StatusBadRequest)
		return
	}

	valid, err := s.vatChecker.ValidateVATID(r.Context(), post.VATNumber)
	if err != nil {
		log.Printf("error reported: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validStatus := valid == "true"
	res := vATPostResponse{
		Valid: validStatus,
	}

	output, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func gracefulServerShutdown(srv *http.Server, idleConnsClosed chan struct{}) {
	defer close(idleConnsClosed)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	<-interrupt

	// We received an interrupt signal, shut down.
	if err := srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("Server gracefully shutdown")
}

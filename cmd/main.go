package main

import (
	"github.com/Lumexralph/vat-id-validator/server"
	"github.com/Lumexralph/vat-id-validator/validator"
	"log"
)

func main() {
	vatChecker := validator.NewVATIDValidator()
	s := server.NewServer(vatChecker)
	if err := s.Start(); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}

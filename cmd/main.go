package main

import (
	"log"
	"net/http"

	"github.com/Lumexralph/vat-id-validator/server"
	"github.com/Lumexralph/vat-id-validator/validator"
)

func main() {
	euServiceURL := "https://ec.europa.eu/taxation_customs/vies/services/checkVatService"
	client := &http.Client{}
	vatChecker := validator.NewVATIDValidator(validator.NewEUVIESService(client, euServiceURL))
	s := server.NewServer(vatChecker)

	if err := s.Start(); err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}

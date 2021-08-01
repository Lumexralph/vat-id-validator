package validator

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// EUServiceVATChecker represents the interface to call external EU service.
type EUServiceVATChecker interface {
	CheckVAT(ctx context.Context, countryCode, vatNumber string) (string, error)
}

// checkVATResponse represents the XML response from the service.
type checkVATResponse struct {
	XMLName     xml.Name `xml:"Envelope"`
	CountryCode string   `xml:"Body>checkVatResponse>countryCode"`
	VATNumber   string   `xml:"Body>checkVatResponse>vatNumber"`
	ValidStatus string   `xml:"Body>checkVatResponse>valid"`
}

// euVIESService the component for the EU/VIES service.
type euVIESService struct {
	client *http.Client
	url    string
}

// NewEUVIESService creates a new instance of the service.
func NewEUVIESService(client *http.Client, url string) *euVIESService {
	return &euVIESService{
		client: client,
		url:    url,
	}
}

// CheckVAT makes an XML request to https://ec.europa.eu/taxation_customs/vies/services/checkVatService,
// it returns a string (valid or invalid) and returns an error.
func (e *euVIESService) CheckVAT(ctx context.Context, countryCode, vatNumber string) (string, error) {
	soapAction := "urn:checkVat"
	httpMethod := "POST"
	requestXML := &bytes.Buffer{}
	postData := &CheckVATPost{
		CountryCode: countryCode,
		VATNumber:   vatNumber,
	}

	if err := Process(checkVatRequestTmpl, postData, requestXML); err != nil {
		return "", fmt.Errorf("error creating xml request: %v", err)
	}

	req, err := http.NewRequest(httpMethod, e.url, requestXML)
	if err != nil {
		log.Printf("error on creating request object: %v\n", err)
		return "", err
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	res, err := e.client.Do(req)
	if err != nil {
		log.Printf("error on dispatching request: %v\n", err)
		return "", err
	}

	var newRes checkVATResponse
	err = xml.NewDecoder(res.Body).Decode(&newRes)
	if err != nil {
		log.Printf("error on unmarshalling xml: %v\n", err)
		return "", err
	}

	return newRes.ValidStatus, nil
}

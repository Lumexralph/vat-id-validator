package validator

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// VATIDChecker interface for the actions our validator can perform.
// It also helps separate the server and the validator if we choose
// to separate them into different components or service and also abstract the implementation.
type VATIDChecker interface {
	// ValidateVATID receives a VAT ID and return a boolean if it's valid or not.
	ValidateVATID(vatID string) (valid bool, err error)
}

const GermanVATPrefix = "DE"

// VATIDValidator is the component that handles validation of
// our VAT Numbers, caching already validated numbers in an in-memory cache,
// also interfacing with an external API to validate the numbers.
type VATIDValidator struct {
	InMemoryCache sync.Map
	client *http.Client
	euService *EUVIESService
}

// NewVATIDValidator creates a new instance of VATIDValidator.
func NewVATIDValidator() *VATIDValidator {
	client := &http.Client{}
	return &VATIDValidator{
		client: client,
		euService: NewEUVIESService(client),
	}
}

func (v *VATIDValidator) ValidateVATID(ctx context.Context, vatID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// sanitize the vatID for whitespace
	vatID = strings.ReplaceAll(vatID, " ", "")
	vatID = strings.ToUpper(vatID)

	var santizedVAT string
	if len(vatID) == 11 {
		santizedVAT = vatID[2:] // strip off the country code.
	}

	// check the cache, if found, return it.
	if val, found := v.InMemoryCache.Load(santizedVAT); found {
		return val == "valid", nil
	}

	if valid := germanVATNumber(vatID); !valid {
		return false, nil
	}

	// validate with the EU/VIES SOAP Service.
	checkStatus, err := v.euService.CheckVAT(ctx, GermanVATPrefix, santizedVAT)
	if err != nil {
		return false, err
	}

	// store in cache.
	v.InMemoryCache.Store(santizedVAT, checkStatus)

	return checkStatus == "valid", nil
}

// germanVATNumber checks if the VAT number is a German VAT Number.
// assumption is we have a sanitized number in this format DE999999999 or 999999999
func germanVATNumber(vatID string) (valid bool) {
	if len(vatID) == 9 { // format: 999999999
		// check if they are all integers
		if _, err := strconv.Atoi(vatID); err != nil {
			return
		}

		valid = true
	}

	if len(vatID) == 11 && strings.ToUpper(vatID[:2]) == GermanVATPrefix { // format: DE999999999
		if _, err := strconv.Atoi(vatID[2:]); err != nil {
			return
		}

		valid = true
	}
	return
}

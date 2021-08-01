package validator

import (
	"io"
	"text/template"
)

// checkVatRequest is the xml template that will be sent to the EU/VIES service.
const checkVatRequestTmpl = `
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
						<countryCode>{{.CountryCode}}</countryCode>
						<vatNumber>{{.VATNumber}}</vatNumber>
					</checkVat>
				</Body>
			</Envelope>`

// CheckVATPost will model request data to be infused into the template.
type CheckVATPost struct {
	CountryCode, VATNumber string
}

// Process creates dynamic template with the data.
func Process(tmpl string, data *CheckVATPost, w io.Writer) error {
	t := template.Must(template.New("check-vat-request").Parse(tmpl))

	// send the template to whatever writer is provided.
	return t.Execute(w, data)
}

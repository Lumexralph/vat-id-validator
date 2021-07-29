package validator

import (
	"bytes"
	"context"
	"encoding/xml"
	"log"
	"net/http"
)

type CheckVATResponse struct {
	XMLName     xml.Name `xml:"Envelope"`
	CountryCode string   `xml:"Body>checkVatResponse>countryCode"`
	VATNumber   string   `xml:"Body>checkVatResponse>vatNumber"`
	ValidStatus string   `xml:"Body>checkVatResponse>valid"`
}

type EUVIESService struct {
	client *http.Client
}

func NewEUVIESService(client *http.Client) *EUVIESService {
	return &EUVIESService{
		client: client,
	}
}

// CheckVAT makes an XML request to https://ec.europa.eu/taxation_customs/vies/services/checkVatService,
// it returns a string (valid or invalid) and returns an error.
func (e *EUVIESService) CheckVAT(ctx context.Context, countryCode, vatNumber string) (string, error) {
	//post := CheckVATPost{
	//	CountryCode: countryCode,
	//	VATNumber:   vatNumber,
	//}
	url := "https://ec.europa.eu/taxation_customs/vies/services/checkVatService"
	soapAction := "urn:checkVat"
	httpMethod := "POST"
	payload := []byte(`<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
    <Body>
        <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
            <countryCode>DE</countryCode>
            <vatNumber>302210417</vatNumber>
        </checkVat>
    </Body>
</Envelope>`)

	req, err := http.NewRequest(httpMethod,url,bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return "", err
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	res, err := e.client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return "", err
	}
	var newRes CheckVATResponse
	err = xml.NewDecoder(res.Body).Decode(&newRes)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return "", err
	}

	return newRes.ValidStatus, nil
}

//func main() {
//	xmlFile, err := os.Open("post.xml")
//	if err != nil {
//		// Defines structs to represent the data
//		fmt.Println("Error opening XML file:", err)
//		return
//	}
//	defer xmlFile.Close()
//	xmlData, err := ioutil.ReadAll(xmlFile)
//	if err != nil {
//		fmt.Println("Error reading XML data:", err)
//		return
//	}
//	var response CheckVATResponse
//	xml.Unmarshal(xmlData, &response)
//	fmt.Println(response)
//
//	post := CheckVATPost{
//		CountryCode: "DE",
//		VATNumber:   "302210417",
//	}
//
//	payload, err := xml.Marshal(&post)
//	if err != nil {
//		fmt.Println("Error marshalling to XML:", err)
//		return
//	}
//	fmt.Println(xml.Header + string(payload))
//
//	url := "https://ec.europa.eu/taxation_customs/vies/services/checkVatService"
//	soapAction := "urn:checkVat"
//	httpMethod := "POST"
//	payload = []byte(`<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
//    <Body>
//        <checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
//            <countryCode>DE</countryCode>
//            <vatNumber>302210417</vatNumber>
//        </checkVat>
//    </Body>
//</Envelope>`)
//
//	req, err := http.NewRequest(httpMethod,url,bytes.NewReader(payload))
//	if err != nil {
//		log.Fatal("Error on creating request object. ", err.Error())
//		return
//	}
//	req.Header.Set("Content-type", "text/xml")
//	req.Header.Set("SOAPAction", soapAction)
//
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		log.Fatal("Error on dispatching request. ", err.Error())
//		return
//	}
//	var newRes CheckVATResponse
//	err = xml.NewDecoder(res.Body).Decode(&newRes)
//	if err != nil {
//		log.Fatal("Error on unmarshaling xml. ", err.Error())
//		return
//	}
//	fmt.Println("server-response", newRes)
//}

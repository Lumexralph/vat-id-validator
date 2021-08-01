package validator

import (
	"bytes"
	"testing"
)

func TestProcess(t *testing.T) {
	var wantTmpl = `
			<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
				<Body>
					<checkVat xmlns="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
						<countryCode>DE</countryCode>
						<vatNumber>123456789</vatNumber>
					</checkVat>
				</Body>
			</Envelope>`

	type args struct {
		tmpl string
		data *CheckVATPost
	}
	tests := []struct {
		name    string
		args    args
		want   string
		wantErr bool
	}{
		{
			name: "process templates successfully",
			args: args{
				tmpl: checkVatRequestTmpl,
				data: &CheckVATPost{
					CountryCode: GermanCountryCode,
					VATNumber: "123456789",
				},
			},
			want: wantTmpl,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := Process(tt.args.tmpl, tt.args.data, w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.want {
				t.Errorf("Process() gotW = %v, want %v", gotW, tt.want)
			}
		})
	}
}

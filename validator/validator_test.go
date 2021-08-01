package validator

import (
	"context"
	"net/http"
	"testing"
)

func Test_GermanVATNumber(t *testing.T) {
	tests := []struct {
		name      string
		vatNumber      string
		wantValid bool
	}{
		{
			name: "german VAT Number format",
			vatNumber: "DE302210417",
			wantValid: true,
		},
		{
			name: "non-german VAT Number format",
			vatNumber: "NL302210417",
			wantValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if gotValid := GermanVATNumber(tt.vatNumber); gotValid != tt.wantValid {
				t.Errorf("germanVATNumber() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func Test_vATIDValidator_ValidateVATID(t *testing.T) {
	tests := []struct {
		name    string
		vatID    string
		want    string
		wantErr bool
	}{
		{
			name: "valid VAT value for test",
			vatID: "DE302210417",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{}
			v := NewVATIDValidator(NewEUVIESService(client, euServiceTestURL))
			got, err := v.ValidateVATID(context.TODO(), tt.vatID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVATID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateVATID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

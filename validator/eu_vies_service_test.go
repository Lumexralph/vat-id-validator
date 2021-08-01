package validator

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestNewEUVIESService(t *testing.T) {
	client := &http.Client{}
	tests := []struct {
		name string
		want *euVIESService
	}{
		{
			name: "create a new instance of EUVIES service",
			want: &euVIESService{
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEUVIESService(client, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEUVIESService() = %v, want %v", got, tt.want)
			}
		})
	}
}

var euServiceTestURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatTestService"

func Test_euVIESService_CheckVAT(t *testing.T) {
	client := &http.Client{}
	ctx := context.TODO()

	type args struct {
		countryCode string
		vatNumber   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "checks a valid VAT",
			args: args{
				countryCode: "ES",
				vatNumber: "100",
			},
			want: "true",
		},
		{
			name: "unsuccessful VAT check",
			args: args{
				countryCode: "ES",
				vatNumber: "TESTVATNUMBER",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewEUVIESService(client, euServiceTestURL)
			got, err := e.CheckVAT(ctx, tt.args.countryCode, tt.args.vatNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckVAT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckVAT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

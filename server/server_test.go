package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lumexralph/vat-id-validator/validator"
)

type fakeVATChecker struct{}

func (f fakeVATChecker) ValidateVATID(ctx context.Context, vatID string) (valid string, err error) {
	if vatID == "error" {
		log.Println("error reached: ")
		return "false", errors.New(vatID)
	}

	if valid := validator.GermanVATNumber(vatID); valid {
		return "true", nil
	}

	return "false", nil
}

func Test_server_healthHandler(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/health", nil)
	s := server{}

	// create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.healthHandler)
	handler.ServeHTTP(rr, r)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("healthHandler(%v, %+v) for url /health; returned wrong status code: got %v want %v", rr, r, status, http.StatusOK)
	}
}

func Test_server_vatIDHandler(t *testing.T) {
	s := NewServer(&fakeVATChecker{})

	type args struct {
		contentType, method, vatID string
		status                     int
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "check a valid VAT request",
			args: args{
				contentType: "application/json",
				method:      "POST",
				vatID:       `{"vat_number":"DE302210417"}`,
				status:      http.StatusOK,
			},
			want: `{"valid":true}`,
		},
		{
			name: "check a non-german VAT ID request",
			args: args{
				contentType: "application/json",
				method:      "POST",
				vatID:       `{"vat_number":"FM402210417"}`,
				status:      http.StatusOK,
			},
			want: `{"valid":false}`,
		},
		{
			name: "check german VAT ID with wrong method request",
			args: args{
				contentType: "application/json",
				method:      "PUT",
				vatID:       `{"vat_number":"FM402210417"}`,
				status:      http.StatusMethodNotAllowed,
			},
			want: "HTTP method not supported\n",
		},
		{
			name: "check german VAT ID with wrong content-type request",
			args: args{
				contentType: "application/xml",
				method:      "POST",
				vatID:       `{"vat_number":"FM402210417"}`,
				status:      http.StatusBadRequest,
			},
			want: "content-type: only json is supported\n",
		},
		{
			name: "check german VAT ID with empty VAT number request",
			args: args{
				contentType: "application/json",
				method:      "POST",
				vatID:       `{"vat_number":""}`,
				status:      http.StatusBadRequest,
			},
			want: "vat_number not provided\n",
		},
		{
			name: "check VAT ID with error from our system",
			args: args{
				contentType: "application/json",
				method:      "POST",
				vatID:       `{"vat_number":"error"}`,
				status:      http.StatusInternalServerError,
			},
			want: "error\n",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name, tt), func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(tt.args.method, "/vatid/validate", bytes.NewReader([]byte(tt.args.vatID)))
			r.Header.Add("Content-Type", tt.args.contentType)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.vatIDHandler)
			handler.ServeHTTP(rr, r)

			if status := rr.Code; status != tt.args.status {
				t.Errorf("vatIDHandler() for %s /vatid/validate; returned wrong status code: got = %v; want=%v", tt.args.method, status, tt.args.status)
			}

			if rr.Body.String() != tt.want {
				t.Errorf("vatIDHandler() for %s /vatid/validate; returned wrong data: got=%s; want=%v", tt.args.method, rr.Body.String(), tt.want)
			}
		})
	}
}

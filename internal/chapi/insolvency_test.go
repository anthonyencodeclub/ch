package chapi_test

import (
	"context"
	"net/http"
	"testing"
)

func TestGetInsolvency_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/12345678/insolvency" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"status": "administration",
			"cases": [
				{
					"number": 1,
					"type": "compulsory-liquidation",
					"dates": [
						{"type": "wound-up-on", "date": "2023-06-15"}
					],
					"practitioners": [
						{
							"name": "Mr A Practitioner",
							"role": "practitioner",
							"address": {
								"address_line_1": "123 Insolvency Lane",
								"locality": "London",
								"postal_code": "EC1A 1BB"
							}
						}
					]
				}
			]
		}`))
	})

	result, err := client.GetInsolvency(context.Background(), "12345678")
	if err != nil {
		t.Fatalf("GetInsolvency() error: %v", err)
	}
	if result.Status != "administration" {
		t.Errorf("Status = %q, want %q", result.Status, "administration")
	}
	if len(result.Cases) != 1 {
		t.Fatalf("Cases count = %d, want 1", len(result.Cases))
	}
	if result.Cases[0].Type != "compulsory-liquidation" {
		t.Errorf("Cases[0].Type = %q, want %q", result.Cases[0].Type, "compulsory-liquidation")
	}
	if len(result.Cases[0].Practitioners) != 1 {
		t.Fatalf("Practitioners count = %d, want 1", len(result.Cases[0].Practitioners))
	}
	if result.Cases[0].Practitioners[0].Name != "Mr A Practitioner" {
		t.Errorf("Practitioner name = %q", result.Cases[0].Practitioners[0].Name)
	}
}

func TestGetInsolvency_NoInsolvency(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"type":"not-found"}]}`))
	})

	_, err := client.GetInsolvency(context.Background(), "00445790")
	if err == nil {
		t.Fatal("expected error for company with no insolvency")
	}
}

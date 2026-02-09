package chapi_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/chapi"
)

func testFilingServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *chapi.FilingClient) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := chapi.NewFilingClientWithBaseURL("test-access-token", srv.URL)
	return srv, client
}

func TestFilingClient_BearerAuth(t *testing.T) {
	var gotAuth string
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	client.GetTransaction(context.Background(), "txn-123")

	if gotAuth != "Bearer test-access-token" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer test-access-token")
	}
}

func TestCreateTransaction_Success(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/transactions" {
			t.Errorf("path = %s, want /transactions", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		json.Unmarshal(body, &payload)
		if payload["company_number"] != "00445790" {
			t.Errorf("company_number = %q, want %q", payload["company_number"], "00445790")
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": "txn-abc123",
			"company_number": "00445790",
			"description": "Change of registered office address",
			"status": "open"
		}`))
	})

	txn, err := client.CreateTransaction(context.Background(), "00445790", "Change of registered office address")
	if err != nil {
		t.Fatalf("CreateTransaction() error: %v", err)
	}
	if txn.ID != "txn-abc123" {
		t.Errorf("ID = %q, want %q", txn.ID, "txn-abc123")
	}
	if txn.Status != "open" {
		t.Errorf("Status = %q, want %q", txn.Status, "open")
	}
}

func TestGetTransaction_Success(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transactions/txn-123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`{"id": "txn-123", "status": "open", "company_number": "00445790"}`))
	})

	txn, err := client.GetTransaction(context.Background(), "txn-123")
	if err != nil {
		t.Fatalf("GetTransaction() error: %v", err)
	}
	if txn.ID != "txn-123" {
		t.Errorf("ID = %q, want %q", txn.ID, "txn-123")
	}
}

func TestCloseTransaction_Success(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.Write([]byte(`{"id": "txn-123", "status": "closed"}`))
	})

	txn, err := client.CloseTransaction(context.Background(), "txn-123")
	if err != nil {
		t.Fatalf("CloseTransaction() error: %v", err)
	}
	if txn.Status != "closed" {
		t.Errorf("Status = %q, want %q", txn.Status, "closed")
	}
}

func TestFileRegisteredOfficeAddress_Success(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/transactions/txn-123/registered-office-address" {
			t.Errorf("path = %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var addr chapi.RegisteredOfficeAddressFiling
		json.Unmarshal(body, &addr)
		if addr.AddressLine1 != "123 New Street" {
			t.Errorf("AddressLine1 = %q", addr.AddressLine1)
		}
		if addr.PostalCode != "SW1A 1AA" {
			t.Errorf("PostalCode = %q", addr.PostalCode)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	})

	addr := chapi.RegisteredOfficeAddressFiling{
		AddressLine1: "123 New Street",
		Locality:     "London",
		PostalCode:   "SW1A 1AA",
		Country:      "United Kingdom",
	}
	err := client.FileRegisteredOfficeAddress(context.Background(), "txn-123", addr)
	if err != nil {
		t.Fatalf("FileRegisteredOfficeAddress() error: %v", err)
	}
}

func TestFileRegisteredEmailAddress_Success(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/transactions/txn-123/registered-email-address" {
			t.Errorf("path = %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var email chapi.RegisteredEmailAddressFiling
		json.Unmarshal(body, &email)
		if email.RegisteredEmailAddress != "info@example.com" {
			t.Errorf("email = %q", email.RegisteredEmailAddress)
		}

		w.WriteHeader(http.StatusCreated)
	})

	email := chapi.RegisteredEmailAddressFiling{
		RegisteredEmailAddress: "info@example.com",
	}
	err := client.FileRegisteredEmailAddress(context.Background(), "txn-123", email)
	if err != nil {
		t.Fatalf("FileRegisteredEmailAddress() error: %v", err)
	}
}

func TestGetAddressValidation_Valid(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transactions/txn-123/registered-office-address/validation-status" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Write([]byte(`{"is_valid": true}`))
	})

	status, err := client.GetAddressValidation(context.Background(), "txn-123")
	if err != nil {
		t.Fatalf("GetAddressValidation() error: %v", err)
	}
	if !status.Valid {
		t.Error("expected valid=true")
	}
}

func TestGetAddressValidation_Invalid(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"is_valid": false, "errors": ["postal_code is required"]}`))
	})

	status, err := client.GetAddressValidation(context.Background(), "txn-123")
	if err != nil {
		t.Fatalf("GetAddressValidation() error: %v", err)
	}
	if status.Valid {
		t.Error("expected valid=false")
	}
	if len(status.Errors) != 1 {
		t.Fatalf("errors count = %d, want 1", len(status.Errors))
	}
}

func TestFilingClient_APIError(t *testing.T) {
	_, client := testFilingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden"}`))
	})

	_, err := client.CreateTransaction(context.Background(), "00445790", "test")
	if err == nil {
		t.Fatal("expected error for 403 response")
	}
}

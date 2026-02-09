package chapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthonyencodeclub/ch/internal/chapi"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *chapi.Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := chapi.NewWithBaseURL("test-api-key", srv.URL)
	return srv, client
}

func TestClient_AuthHeader(t *testing.T) {
	var gotAuth string
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	client.GetCompany(context.Background(), "12345678")

	if gotAuth == "" {
		t.Fatal("expected Authorization header to be set")
	}
}

func TestClient_APIError(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"type":"not-found"}]}`))
	})

	_, err := client.GetCompany(context.Background(), "99999999")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}

	apiErr, ok := err.(*chapi.APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &chapi.APIError{StatusCode: 401, Body: "unauthorized"}
	got := err.Error()
	if got == "" {
		t.Fatal("Error() returned empty string")
	}
	if !contains(got, "401") {
		t.Errorf("Error() = %q, should contain status code", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

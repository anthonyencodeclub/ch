package chapi_test

import (
	"context"
	"net/http"
	"testing"
)

func TestGetCompany_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"company_name": "TESCO PLC",
			"company_number": "00445790",
			"company_status": "active",
			"type": "plc",
			"date_of_creation": "1947-11-27",
			"jurisdiction": "england-wales",
			"registered_office_address": {
				"address_line_1": "Tesco House",
				"locality": "Welwyn Garden City",
				"postal_code": "AL7 1GA"
			},
			"sic_codes": ["47110"]
		}`))
	})

	profile, err := client.GetCompany(context.Background(), "00445790")
	if err != nil {
		t.Fatalf("GetCompany() error: %v", err)
	}
	if profile.CompanyName != "TESCO PLC" {
		t.Errorf("CompanyName = %q, want %q", profile.CompanyName, "TESCO PLC")
	}
	if profile.CompanyNumber != "00445790" {
		t.Errorf("CompanyNumber = %q, want %q", profile.CompanyNumber, "00445790")
	}
	if profile.CompanyStatus != "active" {
		t.Errorf("CompanyStatus = %q, want %q", profile.CompanyStatus, "active")
	}
	if profile.Type != "plc" {
		t.Errorf("Type = %q, want %q", profile.Type, "plc")
	}
	if profile.RegisteredOffice.PostalCode != "AL7 1GA" {
		t.Errorf("PostalCode = %q, want %q", profile.RegisteredOffice.PostalCode, "AL7 1GA")
	}
	if len(profile.SICCodes) != 1 || profile.SICCodes[0] != "47110" {
		t.Errorf("SICCodes = %v, want [47110]", profile.SICCodes)
	}
}

func TestGetCompany_NotFound(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"type":"not-found"}]}`))
	})

	_, err := client.GetCompany(context.Background(), "99999999")
	if err == nil {
		t.Fatal("expected error for non-existent company")
	}
}

func TestGetRegisteredOffice_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790/registered-office-address" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"address_line_1": "Tesco House",
			"address_line_2": "Shire Park",
			"locality": "Welwyn Garden City",
			"region": "Hertfordshire",
			"postal_code": "AL7 1GA",
			"country": "United Kingdom"
		}`))
	})

	addr, err := client.GetRegisteredOffice(context.Background(), "00445790")
	if err != nil {
		t.Fatalf("GetRegisteredOffice() error: %v", err)
	}
	if addr.AddressLine1 != "Tesco House" {
		t.Errorf("AddressLine1 = %q, want %q", addr.AddressLine1, "Tesco House")
	}
	if addr.PostalCode != "AL7 1GA" {
		t.Errorf("PostalCode = %q, want %q", addr.PostalCode, "AL7 1GA")
	}
	if addr.Country != "United Kingdom" {
		t.Errorf("Country = %q, want %q", addr.Country, "United Kingdom")
	}
}

func TestSearchCompanies_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/companies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query().Get("q")
		if q != "tesco" {
			t.Errorf("q = %q, want %q", q, "tesco")
		}
		w.Write([]byte(`{
			"total_results": 42,
			"items": [
				{"company_name": "TESCO PLC", "company_number": "00445790", "company_status": "active"}
			],
			"start_index": 0,
			"items_per_page": 20
		}`))
	})

	result, err := client.SearchCompanies(context.Background(), "tesco", 20, 0)
	if err != nil {
		t.Fatalf("SearchCompanies() error: %v", err)
	}
	if result.TotalResults != 42 {
		t.Errorf("TotalResults = %d, want %d", result.TotalResults, 42)
	}
	if len(result.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(result.Items))
	}
	if result.Items[0].CompanyName != "TESCO PLC" {
		t.Errorf("Items[0].CompanyName = %q, want %q", result.Items[0].CompanyName, "TESCO PLC")
	}
}

func TestSearchCompanies_Pagination(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		ipp := r.URL.Query().Get("items_per_page")
		si := r.URL.Query().Get("start_index")
		if ipp != "5" {
			t.Errorf("items_per_page = %q, want %q", ipp, "5")
		}
		if si != "10" {
			t.Errorf("start_index = %q, want %q", si, "10")
		}
		w.Write([]byte(`{"total_results": 100, "items": [], "start_index": 10, "items_per_page": 5}`))
	})

	_, err := client.SearchCompanies(context.Background(), "test", 5, 10)
	if err != nil {
		t.Fatalf("SearchCompanies() error: %v", err)
	}
}

package chapi_test

import (
	"context"
	"net/http"
	"testing"
)

func TestListOfficers_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790/officers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"total_results": 3,
			"active_count": 2,
			"resigned_count": 1,
			"items": [
				{
					"name": "SMITH, John",
					"officer_role": "director",
					"appointed_on": "2020-01-15"
				},
				{
					"name": "DOE, Jane",
					"officer_role": "secretary",
					"appointed_on": "2019-06-01",
					"resigned_on": "2022-03-15"
				}
			]
		}`))
	})

	result, err := client.ListOfficers(context.Background(), "00445790", 50, 0)
	if err != nil {
		t.Fatalf("ListOfficers() error: %v", err)
	}
	if result.ActiveCount != 2 {
		t.Errorf("ActiveCount = %d, want %d", result.ActiveCount, 2)
	}
	if result.ResignedCount != 1 {
		t.Errorf("ResignedCount = %d, want %d", result.ResignedCount, 1)
	}
	if len(result.Items) != 2 {
		t.Fatalf("Items count = %d, want 2", len(result.Items))
	}
	if result.Items[0].Name != "SMITH, John" {
		t.Errorf("Items[0].Name = %q, want %q", result.Items[0].Name, "SMITH, John")
	}
	if result.Items[1].ResignedOn != "2022-03-15" {
		t.Errorf("Items[1].ResignedOn = %q, want %q", result.Items[1].ResignedOn, "2022-03-15")
	}
}

func TestSearchOfficers_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/officers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query().Get("q")
		if q != "john smith" {
			t.Errorf("q = %q, want %q", q, "john smith")
		}
		w.Write([]byte(`{
			"total_results": 5,
			"items": [
				{"name": "SMITH, John", "officer_role": "director"}
			]
		}`))
	})

	result, err := client.SearchOfficers(context.Background(), "john smith", 20, 0)
	if err != nil {
		t.Fatalf("SearchOfficers() error: %v", err)
	}
	if result.TotalResults != 5 {
		t.Errorf("TotalResults = %d, want %d", result.TotalResults, 5)
	}
}

package chapi_test

import (
	"context"
	"net/http"
	"testing"
)

func TestListPSCs_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790/persons-with-significant-control" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"total_results": 2,
			"active_count": 1,
			"ceased_count": 1,
			"items": [
				{
					"name": "Mr John Smith",
					"kind": "individual-person-with-significant-control",
					"natures_of_control": ["ownership-of-shares-75-to-100-percent"],
					"notified_on": "2016-04-06",
					"nationality": "British",
					"country_of_residence": "England"
				},
				{
					"name": "Mrs Jane Doe",
					"kind": "individual-person-with-significant-control",
					"natures_of_control": ["ownership-of-shares-25-to-50-percent"],
					"notified_on": "2016-04-06",
					"ceased_on": "2023-01-15"
				}
			]
		}`))
	})

	result, err := client.ListPSCs(context.Background(), "00445790", 25, 0)
	if err != nil {
		t.Fatalf("ListPSCs() error: %v", err)
	}
	if result.ActiveCount != 1 {
		t.Errorf("ActiveCount = %d, want %d", result.ActiveCount, 1)
	}
	if result.CeasedCount != 1 {
		t.Errorf("CeasedCount = %d, want %d", result.CeasedCount, 1)
	}
	if len(result.Items) != 2 {
		t.Fatalf("Items count = %d, want 2", len(result.Items))
	}
	if result.Items[0].Nationality != "British" {
		t.Errorf("Items[0].Nationality = %q, want %q", result.Items[0].Nationality, "British")
	}
	if result.Items[1].CeasedOn != "2023-01-15" {
		t.Errorf("Items[1].CeasedOn = %q, want %q", result.Items[1].CeasedOn, "2023-01-15")
	}
}

func TestListPSCs_Empty(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"total_results": 0, "active_count": 0, "ceased_count": 0, "items": []}`))
	})

	result, err := client.ListPSCs(context.Background(), "00445790", 25, 0)
	if err != nil {
		t.Fatalf("ListPSCs() error: %v", err)
	}
	if len(result.Items) != 0 {
		t.Errorf("Items count = %d, want 0", len(result.Items))
	}
}

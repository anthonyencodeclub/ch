package chapi_test

import (
	"context"
	"net/http"
	"testing"
)

func TestListFilingHistory_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790/filing-history" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"total_count": 150,
			"items": [
				{
					"transaction_id": "abc123",
					"category": "accounts",
					"type": "AA",
					"description": "Annual accounts",
					"date": "2024-03-15"
				}
			]
		}`))
	})

	result, err := client.ListFilingHistory(context.Background(), "00445790", "", 25, 0)
	if err != nil {
		t.Fatalf("ListFilingHistory() error: %v", err)
	}
	if result.TotalCount != 150 {
		t.Errorf("TotalCount = %d, want %d", result.TotalCount, 150)
	}
	if len(result.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(result.Items))
	}
	if result.Items[0].Category != "accounts" {
		t.Errorf("Items[0].Category = %q, want %q", result.Items[0].Category, "accounts")
	}
}

func TestListFilingHistory_WithCategory(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		cat := r.URL.Query().Get("category")
		if cat != "accounts" {
			t.Errorf("category = %q, want %q", cat, "accounts")
		}
		w.Write([]byte(`{"total_count": 10, "items": []}`))
	})

	_, err := client.ListFilingHistory(context.Background(), "00445790", "accounts", 25, 0)
	if err != nil {
		t.Fatalf("ListFilingHistory() error: %v", err)
	}
}

func TestGetFilingHistoryItem_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790/filing-history/abc123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"transaction_id": "abc123",
			"category": "confirmation-statement",
			"type": "CS01",
			"description": "Confirmation statement",
			"date": "2024-06-01"
		}`))
	})

	item, err := client.GetFilingHistoryItem(context.Background(), "00445790", "abc123")
	if err != nil {
		t.Fatalf("GetFilingHistoryItem() error: %v", err)
	}
	if item.TransactionID != "abc123" {
		t.Errorf("TransactionID = %q, want %q", item.TransactionID, "abc123")
	}
	if item.Type != "CS01" {
		t.Errorf("Type = %q, want %q", item.Type, "CS01")
	}
}

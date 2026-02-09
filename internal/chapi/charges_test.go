package chapi_test

import (
	"context"
	"net/http"
	"testing"
)

func TestListCharges_Success(t *testing.T) {
	_, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company/00445790/charges" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{
			"total_count": 3,
			"satisfied_count": 1,
			"part_satisfied_count": 0,
			"unfiltered_count": 3,
			"items": [
				{
					"charge_code": "00445790001",
					"status": "outstanding",
					"delivered_on": "2020-05-10"
				},
				{
					"charge_code": "00445790002",
					"status": "satisfied",
					"delivered_on": "2018-03-22",
					"satisfied_on": "2023-01-10"
				}
			]
		}`))
	})

	result, err := client.ListCharges(context.Background(), "00445790", 25, 0)
	if err != nil {
		t.Fatalf("ListCharges() error: %v", err)
	}
	if result.TotalCount != 3 {
		t.Errorf("TotalCount = %d, want %d", result.TotalCount, 3)
	}
	if result.SatisfiedCount != 1 {
		t.Errorf("SatisfiedCount = %d, want %d", result.SatisfiedCount, 1)
	}
	if len(result.Items) != 2 {
		t.Fatalf("Items count = %d, want 2", len(result.Items))
	}
	if result.Items[0].Status != "outstanding" {
		t.Errorf("Items[0].Status = %q, want %q", result.Items[0].Status, "outstanding")
	}
	if result.Items[1].SatisfiedOn != "2023-01-10" {
		t.Errorf("Items[1].SatisfiedOn = %q, want %q", result.Items[1].SatisfiedOn, "2023-01-10")
	}
}

package chapi

import (
	"context"
	"fmt"
	"net/url"
)

// Charge represents a company charge (mortgage/security).
type Charge struct {
	ChargeCode          string `json:"charge_code"`
	Classification      map[string]string `json:"classification"`
	Status              string `json:"status"`
	DeliveredOn         string `json:"delivered_on"`
	CreatedOn           string `json:"created_on,omitempty"`
	SatisfiedOn         string `json:"satisfied_on,omitempty"`
	PersonsEntitled     []map[string]string `json:"persons_entitled,omitempty"`
	Particulars         map[string]any `json:"particulars,omitempty"`
	Links               map[string]string `json:"links,omitempty"`
}

// ChargeList holds a list of charges.
type ChargeList struct {
	TotalCount    int      `json:"total_count"`
	Items         []Charge `json:"items"`
	PartSatisfiedCount int `json:"part_satisfied_count"`
	SatisfiedCount     int `json:"satisfied_count"`
	UnfilteredCount    int `json:"unfiltered_count"`
}

// ListCharges lists charges for a company.
func (c *Client) ListCharges(ctx context.Context, companyNumber string, itemsPerPage, startIndex int) (*ChargeList, error) {
	params := url.Values{}
	if itemsPerPage > 0 {
		params.Set("items_per_page", fmt.Sprintf("%d", itemsPerPage))
	}
	if startIndex > 0 {
		params.Set("start_index", fmt.Sprintf("%d", startIndex))
	}

	var result ChargeList
	if err := c.get(ctx, fmt.Sprintf("/company/%s/charges", companyNumber), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

package chapi

import (
	"context"
	"fmt"
	"net/url"
)

// PSC represents a person with significant control.
type PSC struct {
	Name              string   `json:"name"`
	Kind              string   `json:"kind"`
	NaturesOfControl  []string `json:"natures_of_control"`
	NotifiedOn        string   `json:"notified_on"`
	CeasedOn          string   `json:"ceased_on,omitempty"`
	Nationality       string   `json:"nationality,omitempty"`
	CountryOfResidence string  `json:"country_of_residence,omitempty"`
	Address           RegisteredOffice `json:"address"`
	Links             map[string]string `json:"links,omitempty"`
}

// PSCList holds a list of PSCs.
type PSCList struct {
	TotalResults int   `json:"total_results"`
	ActiveCount  int   `json:"active_count"`
	CeasedCount  int   `json:"ceased_count"`
	Items        []PSC `json:"items"`
	StartIndex   int   `json:"start_index"`
	ItemsPerPage int   `json:"items_per_page"`
}

// ListPSCs lists persons with significant control for a company.
func (c *Client) ListPSCs(ctx context.Context, companyNumber string, itemsPerPage, startIndex int) (*PSCList, error) {
	params := url.Values{}
	if itemsPerPage > 0 {
		params.Set("items_per_page", fmt.Sprintf("%d", itemsPerPage))
	}
	if startIndex > 0 {
		params.Set("start_index", fmt.Sprintf("%d", startIndex))
	}

	var result PSCList
	if err := c.get(ctx, fmt.Sprintf("/company/%s/persons-with-significant-control", companyNumber), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

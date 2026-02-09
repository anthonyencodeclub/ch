package chapi

import (
	"context"
	"fmt"
	"net/url"
)

// Officer represents a company officer.
type Officer struct {
	Name        string `json:"name"`
	OfficerRole string `json:"officer_role"`
	AppointedOn string `json:"appointed_on"`
	ResignedOn  string `json:"resigned_on,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Occupation  string `json:"occupation,omitempty"`
	CountryOfResidence string `json:"country_of_residence,omitempty"`
	Address     RegisteredOffice `json:"address"`
	Links       map[string]string `json:"links,omitempty"`
}

// OfficerList holds a list of officers.
type OfficerList struct {
	TotalResults   int       `json:"total_results"`
	ActiveCount    int       `json:"active_count"`
	ResignedCount  int       `json:"resigned_count"`
	Items          []Officer `json:"items"`
	StartIndex     int       `json:"start_index"`
	ItemsPerPage   int       `json:"items_per_page"`
}

// ListOfficers lists officers for a company.
func (c *Client) ListOfficers(ctx context.Context, companyNumber string, itemsPerPage, startIndex int) (*OfficerList, error) {
	params := url.Values{}
	if itemsPerPage > 0 {
		params.Set("items_per_page", fmt.Sprintf("%d", itemsPerPage))
	}
	if startIndex > 0 {
		params.Set("start_index", fmt.Sprintf("%d", startIndex))
	}

	var result OfficerList
	if err := c.get(ctx, fmt.Sprintf("/company/%s/officers", companyNumber), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OfficerSearchResult holds officer search results.
type OfficerSearchResult struct {
	TotalResults int       `json:"total_results"`
	Items        []Officer `json:"items"`
	StartIndex   int       `json:"start_index"`
	ItemsPerPage int       `json:"items_per_page"`
}

// SearchOfficers searches for officers.
func (c *Client) SearchOfficers(ctx context.Context, query string, itemsPerPage, startIndex int) (*OfficerSearchResult, error) {
	params := url.Values{
		"q": {query},
	}
	if itemsPerPage > 0 {
		params.Set("items_per_page", fmt.Sprintf("%d", itemsPerPage))
	}
	if startIndex > 0 {
		params.Set("start_index", fmt.Sprintf("%d", startIndex))
	}

	var result OfficerSearchResult
	if err := c.get(ctx, "/search/officers", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

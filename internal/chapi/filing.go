package chapi

import (
	"context"
	"fmt"
	"net/url"
)

// FilingHistoryItem represents a single filing.
type FilingHistoryItem struct {
	TransactionID string `json:"transaction_id"`
	Category      string `json:"category"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	Barcode       string `json:"barcode,omitempty"`
	Links         map[string]string `json:"links,omitempty"`
}

// FilingHistoryList holds filing history results.
type FilingHistoryList struct {
	TotalCount   int                 `json:"total_count"`
	Items        []FilingHistoryItem `json:"items"`
	StartIndex   int                 `json:"start_index"`
	ItemsPerPage int                 `json:"items_per_page"`
}

// ListFilingHistory lists filing history for a company.
func (c *Client) ListFilingHistory(ctx context.Context, companyNumber string, category string, itemsPerPage, startIndex int) (*FilingHistoryList, error) {
	params := url.Values{}
	if category != "" {
		params.Set("category", category)
	}
	if itemsPerPage > 0 {
		params.Set("items_per_page", fmt.Sprintf("%d", itemsPerPage))
	}
	if startIndex > 0 {
		params.Set("start_index", fmt.Sprintf("%d", startIndex))
	}

	var result FilingHistoryList
	if err := c.get(ctx, fmt.Sprintf("/company/%s/filing-history", companyNumber), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFilingHistoryItem retrieves a single filing history item.
func (c *Client) GetFilingHistoryItem(ctx context.Context, companyNumber, transactionID string) (*FilingHistoryItem, error) {
	var item FilingHistoryItem
	if err := c.get(ctx, fmt.Sprintf("/company/%s/filing-history/%s", companyNumber, transactionID), nil, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

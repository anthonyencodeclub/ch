package chapi

import (
	"context"
	"fmt"
)

// InsolvencyCase represents an insolvency case.
type InsolvencyCase struct {
	Number       int    `json:"number"`
	Type         string `json:"type"`
	Dates        []InsolvencyDate `json:"dates,omitempty"`
	Practitioners []Practitioner  `json:"practitioners,omitempty"`
}

// InsolvencyDate represents a date within an insolvency case.
type InsolvencyDate struct {
	Type string `json:"type"`
	Date string `json:"date"`
}

// Practitioner represents an insolvency practitioner.
type Practitioner struct {
	Name    string           `json:"name"`
	Role    string           `json:"role"`
	Address RegisteredOffice `json:"address"`
}

// InsolvencyResponse holds insolvency data for a company.
type InsolvencyResponse struct {
	Status string           `json:"status"`
	Cases  []InsolvencyCase `json:"cases"`
}

// GetInsolvency retrieves insolvency information for a company.
func (c *Client) GetInsolvency(ctx context.Context, companyNumber string) (*InsolvencyResponse, error) {
	var result InsolvencyResponse
	if err := c.get(ctx, fmt.Sprintf("/company/%s/insolvency", companyNumber), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

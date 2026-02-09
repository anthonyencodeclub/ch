package chapi

import (
	"context"
	"fmt"
	"net/url"
)

// CompanyProfile represents a company's profile.
type CompanyProfile struct {
	CompanyName        string          `json:"company_name"`
	CompanyNumber      string          `json:"company_number"`
	CompanyStatus      string          `json:"company_status"`
	Type               string          `json:"type"`
	DateOfCreation     string          `json:"date_of_creation"`
	DateOfCessation    string          `json:"date_of_cessation,omitempty"`
	Jurisdiction       string          `json:"jurisdiction"`
	RegisteredOffice   RegisteredOffice `json:"registered_office_address"`
	SICCodes           []string        `json:"sic_codes"`
	HasCharges         bool            `json:"has_charges"`
	HasInsolvencyHistory bool          `json:"has_insolvency_history"`
	Accounts           *Accounts       `json:"accounts,omitempty"`
	ConfirmationStatement *ConfirmationStatement `json:"confirmation_statement,omitempty"`
	Links              map[string]string `json:"links,omitempty"`
}

// RegisteredOffice represents a company's registered address.
type RegisteredOffice struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	Locality     string `json:"locality"`
	Region       string `json:"region,omitempty"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country,omitempty"`
}

// Accounts holds accounting reference date info.
type Accounts struct {
	NextDue          string `json:"next_due,omitempty"`
	LastAccounts     *LastAccounts `json:"last_accounts,omitempty"`
	AccountingReferenceDate *AccountingReferenceDate `json:"accounting_reference_date,omitempty"`
}

// LastAccounts holds the last filed accounts info.
type LastAccounts struct {
	MadeUpTo string `json:"made_up_to,omitempty"`
	Type     string `json:"type,omitempty"`
}

// AccountingReferenceDate holds the ARD.
type AccountingReferenceDate struct {
	Day   string `json:"day"`
	Month string `json:"month"`
}

// ConfirmationStatement holds CS info.
type ConfirmationStatement struct {
	NextDue    string `json:"next_due,omitempty"`
	LastMadeUpTo string `json:"last_made_up_to,omitempty"`
}

// GetCompany retrieves a company profile.
func (c *Client) GetCompany(ctx context.Context, companyNumber string) (*CompanyProfile, error) {
	var profile CompanyProfile
	if err := c.get(ctx, fmt.Sprintf("/company/%s", companyNumber), nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetRegisteredOffice retrieves a company's registered office address.
func (c *Client) GetRegisteredOffice(ctx context.Context, companyNumber string) (*RegisteredOffice, error) {
	var addr RegisteredOffice
	if err := c.get(ctx, fmt.Sprintf("/company/%s/registered-office-address", companyNumber), nil, &addr); err != nil {
		return nil, err
	}
	return &addr, nil
}

// SearchResult holds search results.
type SearchResult struct {
	TotalResults int             `json:"total_results"`
	Items        []CompanyProfile `json:"items"`
	StartIndex   int             `json:"start_index"`
	ItemsPerPage int             `json:"items_per_page"`
}

// SearchCompanies searches for companies.
func (c *Client) SearchCompanies(ctx context.Context, query string, itemsPerPage, startIndex int) (*SearchResult, error) {
	params := url.Values{
		"q": {query},
	}
	if itemsPerPage > 0 {
		params.Set("items_per_page", fmt.Sprintf("%d", itemsPerPage))
	}
	if startIndex > 0 {
		params.Set("start_index", fmt.Sprintf("%d", startIndex))
	}

	var result SearchResult
	if err := c.get(ctx, "/search/companies", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

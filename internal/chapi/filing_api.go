package chapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const filingBaseURL = "https://api.company-information.service.gov.uk"

// FilingClient wraps the Companies House API Filing (write) endpoints.
type FilingClient struct {
	accessToken string
	httpClient  *http.Client
	baseURL     string
}

// NewFilingClient creates a new filing API client with an OAuth2 access token.
func NewFilingClient(accessToken string) *FilingClient {
	return &FilingClient{
		accessToken: accessToken,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		baseURL:     filingBaseURL,
	}
}

// NewFilingClientWithBaseURL creates a filing client with a custom base URL (for testing).
func NewFilingClientWithBaseURL(accessToken, base string) *FilingClient {
	return &FilingClient{
		accessToken: accessToken,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		baseURL:     base,
	}
}

// Transaction represents a filing transaction.
type Transaction struct {
	ID          string            `json:"id,omitempty"`
	CompanyNumber string          `json:"company_number"`
	Description string            `json:"description"`
	Reference   string            `json:"reference,omitempty"`
	Status      string            `json:"status,omitempty"`
	Resources   map[string]any    `json:"resources,omitempty"`
	Links       map[string]string `json:"links,omitempty"`
}

// CreateTransaction creates a new filing transaction.
func (c *FilingClient) CreateTransaction(ctx context.Context, companyNumber, description string) (*Transaction, error) {
	payload := Transaction{
		CompanyNumber: companyNumber,
		Description:   description,
	}
	var result Transaction
	if err := c.post(ctx, "/transactions", payload, &result); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	return &result, nil
}

// GetTransaction retrieves a transaction.
func (c *FilingClient) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	var result Transaction
	if err := c.doGet(ctx, fmt.Sprintf("/transactions/%s", transactionID), &result); err != nil {
		return nil, fmt.Errorf("get transaction: %w", err)
	}
	return &result, nil
}

// CloseTransaction submits/closes a transaction for processing.
func (c *FilingClient) CloseTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	payload := map[string]string{"status": "closed"}
	var result Transaction
	if err := c.put(ctx, fmt.Sprintf("/transactions/%s", transactionID), payload, &result); err != nil {
		return nil, fmt.Errorf("close transaction: %w", err)
	}
	return &result, nil
}

// RegisteredOfficeAddressFiling represents a registered office address change.
type RegisteredOfficeAddressFiling struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	Locality     string `json:"locality"`
	Region       string `json:"region,omitempty"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country,omitempty"`
}

// FileRegisteredOfficeAddress creates a registered office address change within a transaction.
func (c *FilingClient) FileRegisteredOfficeAddress(ctx context.Context, transactionID string, addr RegisteredOfficeAddressFiling) error {
	path := fmt.Sprintf("/transactions/%s/registered-office-address", transactionID)
	return c.post(ctx, path, addr, nil)
}

// RegisteredEmailAddressFiling represents a registered email address change.
type RegisteredEmailAddressFiling struct {
	RegisteredEmailAddress string `json:"registered_email_address"`
}

// FileRegisteredEmailAddress creates a registered email address filing within a transaction.
func (c *FilingClient) FileRegisteredEmailAddress(ctx context.Context, transactionID string, email RegisteredEmailAddressFiling) error {
	path := fmt.Sprintf("/transactions/%s/registered-email-address", transactionID)
	return c.post(ctx, path, email, nil)
}

// ValidationStatus represents the validation result for a filing.
type ValidationStatus struct {
	Valid  bool     `json:"is_valid"`
	Errors []string `json:"errors,omitempty"`
}

// GetAddressValidation checks validation status for a registered office address filing.
func (c *FilingClient) GetAddressValidation(ctx context.Context, transactionID string) (*ValidationStatus, error) {
	var result ValidationStatus
	path := fmt.Sprintf("/transactions/%s/registered-office-address/validation-status", transactionID)
	if err := c.doGet(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get validation: %w", err)
	}
	return &result, nil
}

func (c *FilingClient) post(ctx context.Context, path string, body any, out any) error {
	return c.doRequest(ctx, http.MethodPost, path, body, out)
}

func (c *FilingClient) put(ctx context.Context, path string, body any, out any) error {
	return c.doRequest(ctx, http.MethodPut, path, body, out)
}

func (c *FilingClient) doGet(ctx context.Context, path string, out any) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, out)
}

func (c *FilingClient) doRequest(ctx context.Context, method, path string, body any, out any) error {
	u := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthonyencodeclub/ch/internal/chapi"
	"github.com/anthonyencodeclub/ch/internal/oauth"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
	"github.com/anthonyencodeclub/ch/internal/ui"
)

// FileCmd contains filing operations that modify company data.
type FileCmd struct {
	Address FileAddressCmd `cmd:"" help:"File a change of registered office address"`
	Email   FileEmailCmd   `cmd:"" help:"File a change of registered email address"`
}

// FileAddressCmd files a change of registered office address.
type FileAddressCmd struct {
	CompanyNumber string `arg:"" optional:"" help:"Company number (uses default if omitted)"`
	AddressLine1  string `required:"" help:"Address line 1"`
	AddressLine2  string `help:"Address line 2"`
	Locality      string `required:"" help:"Town or city"`
	Region        string `help:"County or region"`
	PostalCode    string `required:"" help:"Postal code"`
	Country       string `help:"Country"`
}

func (c *FileAddressCmd) Run(ctx context.Context) error {
	cn, err := resolveCompanyNumber(c.CompanyNumber)
	if err != nil {
		return err
	}
	u := ui.FromContext(ctx)

	accessToken, err := oauth.LoadToken(ctx)
	if err != nil {
		return err
	}

	client := chapi.NewFilingClient(accessToken)

	// Step 1: Create transaction
	if u != nil {
		u.Info(fmt.Sprintf("Creating filing transaction for %s...", cn))
	}
	txn, err := client.CreateTransaction(ctx, cn, "Change of registered office address")
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	// Step 2: File the address change
	addr := chapi.RegisteredOfficeAddressFiling{
		AddressLine1: c.AddressLine1,
		AddressLine2: c.AddressLine2,
		Locality:     c.Locality,
		Region:       c.Region,
		PostalCode:   c.PostalCode,
		Country:      c.Country,
	}
	if err := client.FileRegisteredOfficeAddress(ctx, txn.ID, addr); err != nil {
		return fmt.Errorf("file address: %w", err)
	}

	// Step 3: Validate
	validation, err := client.GetAddressValidation(ctx, txn.ID)
	if err != nil {
		if u != nil {
			u.Warn("Could not check validation status (filing may still proceed)")
		}
	} else if !validation.Valid {
		if u != nil {
			u.Error("Validation failed:")
			for _, e := range validation.Errors {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
		}
		return fmt.Errorf("address validation failed")
	}

	// Step 4: Close/submit the transaction
	if u != nil {
		u.Info("Submitting filing...")
	}
	result, err := client.CloseTransaction(ctx, txn.ID)
	if err != nil {
		return fmt.Errorf("submit transaction: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"transaction_id": result.ID,
			"status":         result.Status,
			"company_number": cn,
			"address":        addr,
		})
	}

	if u != nil {
		u.Success(fmt.Sprintf("Address change filed successfully (transaction: %s)", result.ID))
	}
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Transaction:", result.ID)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Status:", result.Status)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Company:", cn)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "New Address:", c.AddressLine1)
	if c.AddressLine2 != "" {
		fmt.Fprintf(os.Stdout, "%-20s %s\n", "", c.AddressLine2)
	}
	fmt.Fprintf(os.Stdout, "%-20s %s %s\n", "", c.Locality, c.PostalCode)

	return nil
}

// FileEmailCmd files a change of registered email address.
type FileEmailCmd struct {
	CompanyNumber string `arg:"" optional:"" help:"Company number (uses default if omitted)"`
	Email         string `required:"" help:"New registered email address"`
}

func (c *FileEmailCmd) Run(ctx context.Context) error {
	cn, err := resolveCompanyNumber(c.CompanyNumber)
	if err != nil {
		return err
	}
	u := ui.FromContext(ctx)

	accessToken, err := oauth.LoadToken(ctx)
	if err != nil {
		return err
	}

	client := chapi.NewFilingClient(accessToken)

	// Step 1: Create transaction
	if u != nil {
		u.Info(fmt.Sprintf("Creating filing transaction for %s...", cn))
	}
	txn, err := client.CreateTransaction(ctx, cn, "Change of registered email address")
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	// Step 2: File the email change
	email := chapi.RegisteredEmailAddressFiling{
		RegisteredEmailAddress: c.Email,
	}
	if err := client.FileRegisteredEmailAddress(ctx, txn.ID, email); err != nil {
		return fmt.Errorf("file email: %w", err)
	}

	// Step 3: Close/submit
	if u != nil {
		u.Info("Submitting filing...")
	}
	result, err := client.CloseTransaction(ctx, txn.ID)
	if err != nil {
		return fmt.Errorf("submit transaction: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"transaction_id": result.ID,
			"status":         result.Status,
			"company_number": cn,
			"email":          c.Email,
		})
	}

	if u != nil {
		u.Success(fmt.Sprintf("Email change filed successfully (transaction: %s)", result.ID))
	}
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Transaction:", result.ID)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Status:", result.Status)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Company:", cn)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "New Email:", c.Email)

	return nil
}

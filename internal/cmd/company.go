package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthonyencodeclub/ch/internal/chapi"
	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
)

// CompanyCmd retrieves company information.
type CompanyCmd struct {
	Get     CompanyGetCmd     `cmd:"" help:"Get company profile"`
	Address CompanyAddressCmd `cmd:"" help:"Get registered office address"`
}

// CompanyGetCmd retrieves a company profile.
type CompanyGetCmd struct {
	CompanyNumber string `arg:"" optional:"" help:"Company number (uses default if omitted)"`
}

func (c *CompanyGetCmd) Run(ctx context.Context) error {
	cn, err := resolveCompanyNumber(c.CompanyNumber)
	if err != nil {
		return err
	}
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	profile, err := client.GetCompany(ctx, cn)
	if err != nil {
		return fmt.Errorf("get company: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, profile)
	}

	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Company Name:", profile.CompanyName)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Company Number:", profile.CompanyNumber)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Status:", profile.CompanyStatus)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Type:", profile.Type)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Incorporated:", profile.DateOfCreation)
	if profile.DateOfCessation != "" {
		fmt.Fprintf(os.Stdout, "%-20s %s\n", "Ceased:", profile.DateOfCessation)
	}
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Jurisdiction:", profile.Jurisdiction)

	addr := profile.RegisteredOffice
	parts := []string{}
	if addr.AddressLine1 != "" {
		parts = append(parts, addr.AddressLine1)
	}
	if addr.AddressLine2 != "" {
		parts = append(parts, addr.AddressLine2)
	}
	if addr.Locality != "" {
		parts = append(parts, addr.Locality)
	}
	if addr.PostalCode != "" {
		parts = append(parts, addr.PostalCode)
	}
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Address:", strings.Join(parts, ", "))

	if len(profile.SICCodes) > 0 {
		fmt.Fprintf(os.Stdout, "%-20s %s\n", "SIC Codes:", strings.Join(profile.SICCodes, ", "))
	}

	return nil
}

// CompanyAddressCmd retrieves the registered office address.
type CompanyAddressCmd struct {
	CompanyNumber string `arg:"" optional:"" help:"Company number (uses default if omitted)"`
}

func (c *CompanyAddressCmd) Run(ctx context.Context) error {
	cn, err := resolveCompanyNumber(c.CompanyNumber)
	if err != nil {
		return err
	}
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	addr, err := client.GetRegisteredOffice(ctx, cn)
	if err != nil {
		return fmt.Errorf("get address: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, addr)
	}

	parts := []string{}
	if addr.AddressLine1 != "" {
		parts = append(parts, addr.AddressLine1)
	}
	if addr.AddressLine2 != "" {
		parts = append(parts, addr.AddressLine2)
	}
	if addr.Locality != "" {
		parts = append(parts, addr.Locality)
	}
	if addr.Region != "" {
		parts = append(parts, addr.Region)
	}
	if addr.PostalCode != "" {
		parts = append(parts, addr.PostalCode)
	}
	if addr.Country != "" {
		parts = append(parts, addr.Country)
	}
	fmt.Fprintln(os.Stdout, strings.Join(parts, ", "))
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthonyencodeclub/ch/internal/chapi"
	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
)

// SearchCmd searches Companies House data.
type SearchCmd struct {
	Companies SearchCompaniesCmd `cmd:"" help:"Search for companies"`
	Officers  SearchOfficersCmd  `cmd:"" help:"Search for officers"`
}

// SearchCompaniesCmd searches for companies.
type SearchCompaniesCmd struct {
	Query        string `arg:"" help:"Search query"`
	ItemsPerPage int    `help:"Results per page" default:"20"`
	StartIndex   int    `help:"Start index for pagination" default:"0"`
}

func (c *SearchCompaniesCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.SearchCompanies(ctx, c.Query, c.ItemsPerPage, c.StartIndex)
	if err != nil {
		return fmt.Errorf("search companies: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Found %d results:\n\n", result.TotalResults)
	for _, item := range result.Items {
		status := item.CompanyStatus
		if status == "" {
			status = "unknown"
		}
		fmt.Fprintf(os.Stdout, "  %-10s  %-50s  %s\n", item.CompanyNumber, item.CompanyName, status)
	}
	return nil
}

// SearchOfficersCmd searches for officers.
type SearchOfficersCmd struct {
	Query        string `arg:"" help:"Search query"`
	ItemsPerPage int    `help:"Results per page" default:"20"`
	StartIndex   int    `help:"Start index for pagination" default:"0"`
}

func (c *SearchOfficersCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.SearchOfficers(ctx, c.Query, c.ItemsPerPage, c.StartIndex)
	if err != nil {
		return fmt.Errorf("search officers: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Found %d results:\n\n", result.TotalResults)
	for _, item := range result.Items {
		role := item.OfficerRole
		if role == "" {
			role = "unknown"
		}
		fmt.Fprintf(os.Stdout, "  %-40s  %s\n", item.Name, role)
	}
	return nil
}

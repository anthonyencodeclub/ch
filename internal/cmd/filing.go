package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthonyencodeclub/ch/internal/chapi"
	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
)

// FilingCmd retrieves filing history.
type FilingCmd struct {
	List FilingListCmd `cmd:"" help:"List filing history for a company"`
	Get  FilingGetCmd  `cmd:"" help:"Get a specific filing"`
}

// FilingListCmd lists filing history.
type FilingListCmd struct {
	CompanyNumber string `arg:"" help:"Company number"`
	Category      string `help:"Filter by category (e.g. accounts, confirmation-statement)" default:""`
	ItemsPerPage  int    `help:"Results per page" default:"25"`
	StartIndex    int    `help:"Start index for pagination" default:"0"`
}

func (c *FilingListCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.ListFilingHistory(ctx, c.CompanyNumber, c.Category, c.ItemsPerPage, c.StartIndex)
	if err != nil {
		return fmt.Errorf("list filings: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Filing History (%d total):\n\n", result.TotalCount)
	for _, f := range result.Items {
		fmt.Fprintf(os.Stdout, "  %s  %-20s  %s\n", f.Date, f.Category, f.Description)
	}
	return nil
}

// FilingGetCmd retrieves a single filing.
type FilingGetCmd struct {
	CompanyNumber string `arg:"" help:"Company number"`
	TransactionID string `arg:"" help:"Transaction ID"`
}

func (c *FilingGetCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	item, err := client.GetFilingHistoryItem(ctx, c.CompanyNumber, c.TransactionID)
	if err != nil {
		return fmt.Errorf("get filing: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, item)
	}

	fmt.Fprintf(os.Stdout, "%-15s %s\n", "Transaction:", item.TransactionID)
	fmt.Fprintf(os.Stdout, "%-15s %s\n", "Date:", item.Date)
	fmt.Fprintf(os.Stdout, "%-15s %s\n", "Category:", item.Category)
	fmt.Fprintf(os.Stdout, "%-15s %s\n", "Type:", item.Type)
	fmt.Fprintf(os.Stdout, "%-15s %s\n", "Description:", item.Description)
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/anthonyencodeclub/ch/internal/chapi"
	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
)

// ChargesCmd retrieves company charges.
type ChargesCmd struct {
	List ChargesListCmd `cmd:"" help:"List charges for a company"`
}

// ChargesListCmd lists charges for a company.
type ChargesListCmd struct {
	CompanyNumber string `arg:"" optional:"" help:"Company number (uses default if omitted)"`
	ItemsPerPage  int    `help:"Results per page" default:"25"`
	StartIndex    int    `help:"Start index for pagination" default:"0"`
}

func (c *ChargesListCmd) Run(ctx context.Context) error {
	cn, err := resolveCompanyNumber(c.CompanyNumber)
	if err != nil {
		return err
	}
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.ListCharges(ctx, cn, c.ItemsPerPage, c.StartIndex)
	if err != nil {
		return fmt.Errorf("list charges: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Charges (%d total, %d satisfied):\n\n", result.TotalCount, result.SatisfiedCount)
	for _, ch := range result.Items {
		fmt.Fprintf(os.Stdout, "  %-20s  %-15s  delivered %s\n", ch.ChargeCode, ch.Status, ch.DeliveredOn)
	}
	return nil
}

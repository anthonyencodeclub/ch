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

// PSCCmd retrieves persons with significant control.
type PSCCmd struct {
	List PSCListCmd `cmd:"" help:"List persons with significant control"`
}

// PSCListCmd lists PSCs for a company.
type PSCListCmd struct {
	CompanyNumber string `arg:"" help:"Company number"`
	ItemsPerPage  int    `help:"Results per page" default:"25"`
	StartIndex    int    `help:"Start index for pagination" default:"0"`
}

func (c *PSCListCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.ListPSCs(ctx, c.CompanyNumber, c.ItemsPerPage, c.StartIndex)
	if err != nil {
		return fmt.Errorf("list PSCs: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Persons with Significant Control (%d active, %d ceased):\n\n", result.ActiveCount, result.CeasedCount)
	for _, p := range result.Items {
		ceased := ""
		if p.CeasedOn != "" {
			ceased = fmt.Sprintf(" (ceased %s)", p.CeasedOn)
		}
		controls := strings.Join(p.NaturesOfControl, "; ")
		fmt.Fprintf(os.Stdout, "  %-40s  notified %s%s\n", p.Name, p.NotifiedOn, ceased)
		if controls != "" {
			fmt.Fprintf(os.Stdout, "    Controls: %s\n", controls)
		}
	}
	return nil
}

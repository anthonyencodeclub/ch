package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/openclaw/ch/internal/chapi"
	"github.com/openclaw/ch/internal/config"
	"github.com/openclaw/ch/internal/outfmt"
)

// OfficersCmd lists and views company officers.
type OfficersCmd struct {
	List OfficersListCmd `cmd:"" help:"List officers for a company"`
}

// OfficersListCmd lists officers for a company.
type OfficersListCmd struct {
	CompanyNumber string `arg:"" help:"Company number"`
	ItemsPerPage  int    `help:"Results per page" default:"50"`
	StartIndex    int    `help:"Start index for pagination" default:"0"`
}

func (c *OfficersListCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.ListOfficers(ctx, c.CompanyNumber, c.ItemsPerPage, c.StartIndex)
	if err != nil {
		return fmt.Errorf("list officers: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Officers (%d active, %d resigned):\n\n", result.ActiveCount, result.ResignedCount)
	for _, o := range result.Items {
		resigned := ""
		if o.ResignedOn != "" {
			resigned = fmt.Sprintf(" (resigned %s)", o.ResignedOn)
		}
		fmt.Fprintf(os.Stdout, "  %-40s  %-20s  appointed %s%s\n", o.Name, o.OfficerRole, o.AppointedOn, resigned)
	}
	return nil
}

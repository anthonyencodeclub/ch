package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/openclaw/ch/internal/chapi"
	"github.com/openclaw/ch/internal/config"
	"github.com/openclaw/ch/internal/outfmt"
)

// InsolvencyCmd retrieves insolvency information.
type InsolvencyCmd struct {
	Get InsolvencyGetCmd `cmd:"" help:"Get insolvency information for a company"`
}

// InsolvencyGetCmd retrieves insolvency data.
type InsolvencyGetCmd struct {
	CompanyNumber string `arg:"" help:"Company number"`
}

func (c *InsolvencyGetCmd) Run(ctx context.Context) error {
	apiKey, err := config.APIKey()
	if err != nil {
		return err
	}

	client := chapi.New(apiKey)
	result, err := client.GetInsolvency(ctx, c.CompanyNumber)
	if err != nil {
		return fmt.Errorf("get insolvency: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stdout, "Insolvency Status: %s\n\n", result.Status)
	for _, cs := range result.Cases {
		fmt.Fprintf(os.Stdout, "  Case %d (%s)\n", cs.Number, cs.Type)
		for _, d := range cs.Dates {
			fmt.Fprintf(os.Stdout, "    %-20s %s\n", d.Type+":", d.Date)
		}
		for _, p := range cs.Practitioners {
			fmt.Fprintf(os.Stdout, "    Practitioner: %s (%s)\n", p.Name, p.Role)
		}
		fmt.Fprintln(os.Stdout)
	}
	return nil
}

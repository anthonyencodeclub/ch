package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthonyencodeclub/ch/internal/chapi"
	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
	"github.com/anthonyencodeclub/ch/internal/ui"
)

// SetupCmd is the guided company setup flow.
type SetupCmd struct{}

func (c *SetupCmd) Run(ctx context.Context) error {
	u := ui.FromContext(ctx)
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Check for API key
	apiKey, err := config.APIKey()
	if err != nil || apiKey == "" {
		if u != nil {
			u.Info("No API key found. Get one free at https://developer.company-information.service.gov.uk/")
			fmt.Fprint(os.Stderr, "\nEnter your Companies House API key: ")
		}

		input, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("read input: %w", readErr)
		}
		apiKey = strings.TrimSpace(input)
		if apiKey == "" {
			return fmt.Errorf("API key is required")
		}

		cfg, _ := config.ReadConfig()
		cfg.APIKey = apiKey
		if writeErr := config.WriteConfig(cfg); writeErr != nil {
			return fmt.Errorf("save API key: %w", writeErr)
		}
		if u != nil {
			u.Success("API key saved.")
		}
	} else {
		if u != nil {
			u.Success("API key found.")
		}
	}

	// Step 2: Search for or enter company number
	if u != nil {
		fmt.Fprintln(os.Stderr)
		u.Info("Now let's find your company.")
		fmt.Fprint(os.Stderr, "Enter a company name or number: ")
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	query := strings.TrimSpace(input)
	if query == "" {
		return fmt.Errorf("company name or number is required")
	}

	client := chapi.New(apiKey)

	// Try as a company number first (8 digits, possibly with leading zeros)
	companyNumber := query
	profile, err := client.GetCompany(ctx, companyNumber)
	if err != nil {
		// Not a valid company number — try searching
		if u != nil {
			u.Info(fmt.Sprintf("Searching for %q...", query))
		}

		results, searchErr := client.SearchCompanies(ctx, query, 10, 0)
		if searchErr != nil {
			return fmt.Errorf("search companies: %w", searchErr)
		}

		if results.TotalResults == 0 {
			return fmt.Errorf("no companies found matching %q", query)
		}

		// Display results
		fmt.Fprintln(os.Stderr)
		for i, item := range results.Items {
			status := item.CompanyStatus
			if status == "" {
				status = "unknown"
			}
			fmt.Fprintf(os.Stderr, "  [%d] %-10s  %-45s  %s\n", i+1, item.CompanyNumber, item.CompanyName, status)
		}

		fmt.Fprintf(os.Stderr, "\nSelect a company (1-%d): ", len(results.Items))
		selInput, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("read input: %w", readErr)
		}

		var sel int
		if _, scanErr := fmt.Sscanf(strings.TrimSpace(selInput), "%d", &sel); scanErr != nil || sel < 1 || sel > len(results.Items) {
			return fmt.Errorf("invalid selection")
		}

		selected := results.Items[sel-1]
		companyNumber = selected.CompanyNumber

		// Fetch full profile
		profile, err = client.GetCompany(ctx, companyNumber)
		if err != nil {
			return fmt.Errorf("get company: %w", err)
		}
	}

	// Step 3: Save as default company
	cfg, _ := config.ReadConfig()
	cfg.DefaultCompany = profile.CompanyNumber
	cfg.CompanyName = profile.CompanyName
	if writeErr := config.WriteConfig(cfg); writeErr != nil {
		return fmt.Errorf("save config: %w", writeErr)
	}

	// Step 4: Display summary
	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"company_number": profile.CompanyNumber,
			"company_name":   profile.CompanyName,
			"status":         profile.CompanyStatus,
			"type":           profile.Type,
			"incorporated":   profile.DateOfCreation,
			"jurisdiction":   profile.Jurisdiction,
			"address":        formatAddress(profile.RegisteredOffice),
			"configured":     true,
		})
	}

	fmt.Fprintln(os.Stderr)
	if u != nil {
		u.Success(fmt.Sprintf("Company configured: %s (%s)", profile.CompanyName, profile.CompanyNumber))
	}
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Company Name:", profile.CompanyName)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Company Number:", profile.CompanyNumber)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Status:", profile.CompanyStatus)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Type:", profile.Type)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Incorporated:", profile.DateOfCreation)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Jurisdiction:", profile.Jurisdiction)
	fmt.Fprintf(os.Stdout, "%-20s %s\n", "Address:", formatAddress(profile.RegisteredOffice))

	if len(profile.SICCodes) > 0 {
		fmt.Fprintf(os.Stdout, "%-20s %s\n", "SIC Codes:", strings.Join(profile.SICCodes, ", "))
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stderr, "You can now run commands without specifying a company number:")
	fmt.Fprintln(os.Stderr, "  ch company get")
	fmt.Fprintln(os.Stderr, "  ch officers list")
	fmt.Fprintln(os.Stderr, "  ch filing list")

	return nil
}

func formatAddress(addr chapi.RegisteredOffice) string {
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
	return strings.Join(parts, ", ")
}

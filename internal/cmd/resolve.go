package cmd

import (
	"fmt"

	"github.com/anthonyencodeclub/ch/internal/config"
)

// resolveCompanyNumber returns the provided company number, or falls back to the
// default company from config if the argument is empty.
func resolveCompanyNumber(companyNumber string) (string, error) {
	if companyNumber != "" {
		return companyNumber, nil
	}
	cfg, err := config.ReadConfig()
	if err != nil {
		return "", fmt.Errorf("read config: %w", err)
	}
	if cfg.DefaultCompany == "" {
		return "", fmt.Errorf("no company number provided and no default set (run: ch setup)")
	}
	return cfg.DefaultCompany, nil
}

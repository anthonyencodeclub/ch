package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/openclaw/ch/internal/config"
	"github.com/openclaw/ch/internal/outfmt"
	"github.com/openclaw/ch/internal/ui"
)

// AuthCmd manages API key authentication.
type AuthCmd struct {
	SetKey AuthSetKeyCmd `cmd:"" name:"set-key" help:"Store your Companies House API key"`
	Status AuthStatusCmd `cmd:"" help:"Show current auth status"`
}

// AuthSetKeyCmd stores an API key.
type AuthSetKeyCmd struct {
	Key string `arg:"" help:"Your Companies House API key"`
}

func (c *AuthSetKeyCmd) Run(ctx context.Context) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	cfg.APIKey = c.Key
	if err := config.WriteConfig(cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	if u := ui.FromContext(ctx); u != nil {
		u.Success("API key saved.")
	}
	return nil
}

// AuthStatusCmd shows the current auth status.
type AuthStatusCmd struct{}

func (c *AuthStatusCmd) Run(ctx context.Context) error {
	key, err := config.APIKey()
	hasKey := err == nil && key != ""

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"authenticated": hasKey,
		})
	}

	if u := ui.FromContext(ctx); u != nil {
		if hasKey {
			masked := key[:4] + "..." + key[len(key)-4:]
			u.Success(fmt.Sprintf("Authenticated (key: %s)", masked))
		} else {
			u.Warn("Not authenticated. Run: ch auth set-key <YOUR_KEY>")
		}
	}
	return nil
}

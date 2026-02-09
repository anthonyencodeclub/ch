package cmd

import (
	"context"
	"fmt"
	"os"

	"time"

	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/oauth"
	"github.com/anthonyencodeclub/ch/internal/outfmt"
	"github.com/anthonyencodeclub/ch/internal/ui"
)

// AuthCmd manages API key authentication.
type AuthCmd struct {
	SetKey AuthSetKeyCmd `cmd:"" name:"set-key" help:"Store your Companies House API key"`
	Login  AuthLoginCmd  `cmd:"" help:"Login via OAuth2 for filing operations (change address, email)"`
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

// AuthLoginCmd performs OAuth2 login for filing operations.
type AuthLoginCmd struct {
	ClientID     string `help:"OAuth2 client ID" env:"CH_CLIENT_ID"`
	ClientSecret string `help:"OAuth2 client secret" env:"CH_CLIENT_SECRET"`
}

func (c *AuthLoginCmd) Run(ctx context.Context) error {
	u := ui.FromContext(ctx)

	// Load from config if not provided as flags
	cfg, _ := config.ReadConfig()
	clientID := c.ClientID
	clientSecret := c.ClientSecret
	if clientID == "" {
		clientID = cfg.OAuthClientID
	}
	if clientSecret == "" {
		clientSecret = cfg.OAuthClientSecret
	}

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("OAuth2 client ID and secret required.\n\n" +
			"Register an application at https://developer.company-information.service.gov.uk/\n" +
			"Then run: ch auth login --client-id YOUR_ID --client-secret YOUR_SECRET")
	}

	if u != nil {
		u.Info("Opening browser for Companies House login...")
	}

	tok, err := oauth.Login(ctx, clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	// Save everything to config
	cfg.OAuthClientID = clientID
	cfg.OAuthClientSecret = clientSecret
	cfg.OAuthAccessToken = tok.AccessToken
	cfg.OAuthRefreshToken = tok.RefreshToken
	cfg.OAuthTokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second).Format(time.RFC3339)
	if err := config.WriteConfig(cfg); err != nil {
		return fmt.Errorf("save tokens: %w", err)
	}

	if u != nil {
		u.Success("Logged in successfully! You can now use filing commands.")
	}
	return nil
}

// AuthStatusCmd shows the current auth status.
type AuthStatusCmd struct{}

func (c *AuthStatusCmd) Run(ctx context.Context) error {
	key, err := config.APIKey()
	hasKey := err == nil && key != ""

	cfg, _ := config.ReadConfig()
	hasOAuth := cfg.OAuthAccessToken != ""

	if outfmt.IsJSON(ctx) {
		result := map[string]any{
			"api_key_set": hasKey,
			"oauth_login": hasOAuth,
		}
		if hasOAuth && cfg.OAuthTokenExpiry != "" {
			result["oauth_expires"] = cfg.OAuthTokenExpiry
		}
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if u := ui.FromContext(ctx); u != nil {
		if hasKey {
			masked := key[:4] + "..." + key[len(key)-4:]
			u.Success(fmt.Sprintf("API key: %s", masked))
		} else {
			u.Warn("API key: not set (run: ch auth set-key)")
		}

		if hasOAuth {
			expiry := "unknown"
			if cfg.OAuthTokenExpiry != "" {
				if t, err := time.Parse(time.RFC3339, cfg.OAuthTokenExpiry); err == nil {
					if time.Now().Before(t) {
						expiry = fmt.Sprintf("expires %s", t.Format("2006-01-02 15:04"))
					} else {
						expiry = "expired (will auto-refresh)"
					}
				}
			}
			u.Success(fmt.Sprintf("OAuth2 filing: logged in (%s)", expiry))
		} else {
			u.Warn("OAuth2 filing: not logged in (run: ch auth login)")
		}
	}
	return nil
}

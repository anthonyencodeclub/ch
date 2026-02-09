package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anthonyencodeclub/ch/internal/config"
	"github.com/anthonyencodeclub/ch/internal/oauth"
)

func TestRefreshAccessToken_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		r.ParseForm()
		if r.Form.Get("grant_type") != "refresh_token" {
			t.Errorf("grant_type = %q, want %q", r.Form.Get("grant_type"), "refresh_token")
		}
		if r.Form.Get("refresh_token") != "old-refresh-token" {
			t.Errorf("refresh_token = %q", r.Form.Get("refresh_token"))
		}
		if r.Form.Get("client_id") != "test-client-id" {
			t.Errorf("client_id = %q", r.Form.Get("client_id"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_in":    3600,
			"token_type":    "Bearer",
		})
	}))
	defer srv.Close()

	// Override the token URL for testing
	tok, err := oauth.RefreshAccessTokenWithURL(context.Background(), srv.URL, "test-client-id", "test-secret", "old-refresh-token")
	if err != nil {
		t.Fatalf("RefreshAccessToken() error: %v", err)
	}
	if tok.AccessToken != "new-access-token" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "new-access-token")
	}
	if tok.RefreshToken != "new-refresh-token" {
		t.Errorf("RefreshToken = %q, want %q", tok.RefreshToken, "new-refresh-token")
	}
	if tok.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %d, want %d", tok.ExpiresIn, 3600)
	}
}

func TestRefreshAccessToken_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "invalid_grant",
			"error_description": "refresh token expired",
		})
	}))
	defer srv.Close()

	_, err := oauth.RefreshAccessTokenWithURL(context.Background(), srv.URL, "id", "secret", "bad-token")
	if err == nil {
		t.Fatal("expected error for invalid refresh token")
	}
}

func TestLoadToken_NoToken(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	_, err := oauth.LoadToken(context.Background())
	if err == nil {
		t.Fatal("expected error when no token saved")
	}
}

func TestLoadToken_ValidToken(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg := config.File{
		OAuthAccessToken: "valid-token",
		OAuthTokenExpiry: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	tok, err := oauth.LoadToken(context.Background())
	if err != nil {
		t.Fatalf("LoadToken() error: %v", err)
	}
	if tok != "valid-token" {
		t.Errorf("LoadToken() = %q, want %q", tok, "valid-token")
	}
}

func TestLoadToken_ExpiredWithNoRefresh(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("CH_CONFIG_DIR", tmp)

	cfg := config.File{
		OAuthAccessToken: "expired-token",
		OAuthTokenExpiry: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
	}
	if err := config.WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig() error: %v", err)
	}

	_, err := oauth.LoadToken(context.Background())
	if err == nil {
		t.Fatal("expected error for expired token with no refresh token")
	}
}

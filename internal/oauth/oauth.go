package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/anthonyencodeclub/ch/internal/config"
)

const (
	authURL  = "https://identity.company-information.service.gov.uk/oauth2/authorise"
	tokenURL = "https://identity.company-information.service.gov.uk/oauth2/token"
)

// TokenResponse is the response from the token endpoint.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Login performs the OAuth2 authorization code flow using a local callback server.
func Login(ctx context.Context, clientID, clientSecret string) (*TokenResponse, error) {
	// Find a free port for the callback
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	state, err := randomState()
	if err != nil {
		return nil, fmt.Errorf("generate state: %w", err)
	}

	// Build authorization URL
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"scope":         {"https://identity.company-information.service.gov.uk/user/profile.read"},
		"state":         {state},
	}
	authzURL := authURL + "?" + params.Encode()

	// Channel to receive the authorization code
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "invalid state", http.StatusBadRequest)
			errCh <- fmt.Errorf("state mismatch")
			return
		}
		if errParam := r.URL.Query().Get("error"); errParam != "" {
			desc := r.URL.Query().Get("error_description")
			http.Error(w, "Authorization failed: "+desc, http.StatusBadRequest)
			errCh <- fmt.Errorf("authorization error: %s - %s", errParam, desc)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "no code received", http.StatusBadRequest)
			errCh <- fmt.Errorf("no authorization code received")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html><html><body><h2>Authenticated!</h2><p>You can close this tab and return to the terminal.</p></body></html>`)
		codeCh <- code
	})

	srv := &http.Server{Handler: mux}

	go func() {
		if serveErr := srv.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			errCh <- serveErr
		}
	}()

	// Open browser
	openBrowser(authzURL)

	// Wait for callback or timeout
	var code string
	select {
	case code = <-codeCh:
	case err = <-errCh:
		srv.Shutdown(ctx)
		return nil, err
	case <-time.After(5 * time.Minute):
		srv.Shutdown(ctx)
		return nil, fmt.Errorf("timed out waiting for authorization (5 minutes)")
	case <-ctx.Done():
		srv.Shutdown(ctx)
		return nil, ctx.Err()
	}

	srv.Shutdown(ctx)

	// Exchange code for token
	return exchangeCode(ctx, clientID, clientSecret, code, redirectURI)
}

// RefreshAccessToken refreshes an access token using a refresh token.
func RefreshAccessToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	return RefreshAccessTokenWithURL(ctx, tokenURL, clientID, clientSecret, refreshToken)
}

// RefreshAccessTokenWithURL refreshes a token using a custom token endpoint (for testing).
func RefreshAccessTokenWithURL(ctx context.Context, tokenEndpoint, clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	return doTokenRequestURL(ctx, tokenEndpoint, data)
}

// LoadToken returns a valid access token, refreshing if necessary.
func LoadToken(ctx context.Context) (string, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		return "", err
	}
	if cfg.OAuthAccessToken == "" {
		return "", fmt.Errorf("not logged in for filing (run: ch auth login)")
	}

	// Check if token is expired
	if cfg.OAuthTokenExpiry != "" {
		expiry, parseErr := time.Parse(time.RFC3339, cfg.OAuthTokenExpiry)
		if parseErr == nil && time.Now().After(expiry) {
			// Token expired, try to refresh
			if cfg.OAuthRefreshToken == "" || cfg.OAuthClientID == "" || cfg.OAuthClientSecret == "" {
				return "", fmt.Errorf("access token expired, please re-login: ch auth login")
			}
			tok, refreshErr := RefreshAccessToken(ctx, cfg.OAuthClientID, cfg.OAuthClientSecret, cfg.OAuthRefreshToken)
			if refreshErr != nil {
				return "", fmt.Errorf("refresh token: %w (try: ch auth login)", refreshErr)
			}
			cfg.OAuthAccessToken = tok.AccessToken
			cfg.OAuthRefreshToken = tok.RefreshToken
			cfg.OAuthTokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second).Format(time.RFC3339)
			if writeErr := config.WriteConfig(cfg); writeErr != nil {
				return "", fmt.Errorf("save refreshed token: %w", writeErr)
			}
		}
	}

	return cfg.OAuthAccessToken, nil
}

func exchangeCode(ctx context.Context, clientID, clientSecret, code, redirectURI string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
	}
	return doTokenRequest(ctx, data)
}

func doTokenRequest(ctx context.Context, data url.Values) (*TokenResponse, error) {
	return doTokenRequestURL(ctx, tokenURL, data)
}

func doTokenRequestURL(ctx context.Context, endpoint string, data url.Values) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		return nil, fmt.Errorf("token exchange failed (HTTP %d): %s - %s", resp.StatusCode, body.Error, body.Description)
	}

	var tok TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	return &tok, nil
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}

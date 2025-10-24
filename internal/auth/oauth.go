package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var OAuthConfig *oauth2.Config

const (
	TokenFile       = ".drive-cli-token.json"
	CredentialsFile = "oauth.json"
	AuthPort        = 8080 // Fixed port for OAuth callback
)

var (
	authCode     string
	authCodeOnce sync.Once
	authServer   *http.Server
)

func InitOAuth(clientID, clientSecret, redirectURL string) {
	OAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			drive.DriveScope,
		},
		Endpoint: google.Endpoint,
	}
}

// InitOAuthFromFile initializes OAuth config from credentials file
func InitOAuthFromFile() error {
	b, err := os.ReadFile(CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return fmt.Errorf("unable to parse credentials: %v", err)
	}

	OAuthConfig = config
	return nil
}

// GetToken returns a valid token, refreshing if necessary
func GetToken(ctx context.Context) (*oauth2.Token, error) {
	tok, err := tokenFromFile(TokenFile)
	if err != nil {
		return nil, fmt.Errorf("token file not found: %v. Run 'gooDrive auth' to authenticate", err)
	}

	// Check if token is expired and refresh if needed
	if tok.Expiry.Before(time.Now()) {
		log.Println("Token expired, refreshing...")
		newTok, err := OAuthConfig.TokenSource(ctx, tok).Token()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v. Run 'gooDrive auth' to re-authenticate", err)
		}
		tok = newTok
		SaveToken(tok)
		log.Println("Token refreshed successfully")
	}

	return tok, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func SaveToken(token *oauth2.Token) error {
	f, err := os.Create(TokenFile)
	if err != nil {
		return fmt.Errorf("unable to cache token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func AuthURL() string {
	return OAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := OAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange token: %v", err)
	}

	return token, nil
}

// GetTokenFromWeb starts a local server and handles OAuth flow
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Use fixed port
	port := AuthPort

	// Update redirect URL
	config.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)

	// Setup local HTTP server
	mux := http.NewServeMux()

	// Serve the auth template HTML file
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the auth_template.html file
		htmlContent, err := os.ReadFile("auth_template.html")
		if err != nil {
			// Fallback to simple HTML if template not found
			log.Printf("Warning: auth_template.html not found, using fallback: %v", err)
			htmlContent = []byte(`<!DOCTYPE html>
<html>
<head>
    <title>gooDrive - Authentication</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background: #0d1117; color: #c9d1d9; }
        .container { max-width: 600px; margin: 0 auto; }
        h1 { color: #58a6ff; }
    </style>
</head>
<body>
    <div class="container">
        <h1>gooDrive</h1>
        <p>Please complete the authentication process...</p>
    </div>
</body>
</html>`)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlContent)
	})

	// Handle OAuth callback
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "error": "No authorization code received"}`))
			return
		}

		authCodeOnce.Do(func() {
			authCode = code
		})

		// Redirect to success page
		http.Redirect(w, r, "/?success=true", http.StatusSeeOther)

		// Shutdown server after successful auth
		go func() {
			time.Sleep(3 * time.Second)
			if authServer != nil {
				authServer.Shutdown(context.Background())
			}
		}()
	})

	authServer = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: mux,
	}

	// Start server in background
	go func() {
		if err := authServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Generate auth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Google Drive CLI - Authentication Required           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("\n🌐 Opening browser for authentication...\n\n")
	fmt.Printf("If your browser doesn't open automatically, visit:\n%s\n\n", authURL)
	fmt.Printf("🔗 Local server started at: http://localhost:%d\n", port)
	fmt.Println("\n⏳ Waiting for authentication...")

	// Try to open browser automatically
	openBrowser(authURL)

	// Wait for auth code or timeout
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatal("Authentication timeout. Please try again.")
		case <-ticker.C:
			if authCode != "" {
				fmt.Println("\n✅ Authentication successful!")

				tok, err := config.Exchange(context.TODO(), authCode)
				if err != nil {
					log.Fatalf("Unable to retrieve token: %v", err)
				}

				// Cleanup
				if authServer != nil {
					authServer.Shutdown(context.Background())
				}

				return tok
			}
		}
	}
}

// openBrowser tries to open the URL in a browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

// GetHTTPClient returns an authenticated HTTP client
func GetHTTPClient() (*http.Client, error) {
	if OAuthConfig == nil {
		if err := InitOAuthFromFile(); err != nil {
			return nil, err
		}
	}

	tok, err := GetToken(context.Background())
	if err != nil {
		return nil, err
	}

	return OAuthConfig.Client(context.Background(), tok), nil
}

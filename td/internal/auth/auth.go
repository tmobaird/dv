package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/tmobaird/dv/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func GetClient(ctx context.Context) (*http.Client, error) {
	config, err := getConfig()
	if err != nil {
		return &http.Client{}, nil
	}

	tokenFile := tokenCacheFile()
	tok, err := tokenFromFile(tokenFile)

	if err != nil {
		client := getTokenForFirstTime(ctx, config)
		return client, nil
	}

	return config.Client(ctx, tok), nil
}

func getConfig() (*oauth2.Config, error) {
	// Load OAuth 2.0 config from client_secrets.json
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return &oauth2.Config{}, err
	}

	// Set the desired scopes
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)

	if err != nil {
		return &oauth2.Config{}, err
	}

	return config, nil
}

func tokenCacheFile() string {
	tokenCacheDir := filepath.Join(core.BasePath(), "credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, "google.json")
}

func saveToken(path string, token *oauth2.Token) {
	f, _ := os.Create(path)
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func getTokenForFirstTime(ctx context.Context, config *oauth2.Config) *http.Client {
	tokFile := tokenCacheFile()
	codeCh := make(chan string)
	srv := &http.Server{Addr: ":80"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No code in request", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "Authorization complete. You can close this window.")
		codeCh <- code

		go func() {
			_ = srv.Shutdown(context.Background())
		}()
	})

	// Start the server
	go func() {
		_ = srv.ListenAndServe()
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Opening browser to: %s\n", authURL)
	_ = openBrowser(authURL)

	// Wait for code
	code := <-codeCh

	// Exchange code for token
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	saveToken(tokFile, tok)

	return config.Client(ctx, tok)
}

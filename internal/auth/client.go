package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
)

func tokenCacheFile() string {
	usr, _ := user.Current()
	tokenCacheDir := filepath.Join(usr.HomeDir, ".td", "credentials")
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
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
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

func GetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tokFile := tokenCacheFile()
	tok, err := tokenFromFile(tokFile)

	if err != nil {
		getTokenForFirstTime(ctx, config)
	}

	return config.Client(ctx, tok)
}

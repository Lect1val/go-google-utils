package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Lect1val/go-google-utils/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// GenerateTokenInteractive runs an interactive OAuth flow using the provided
// client secret file, then saves the resulting token JSON to the given path.
// Uses only the Gmail send scope for least-privilege access.
func GenerateTokenInteractive(tokenOutputPath string) error {
	// Ensure config is loaded and read client secret from env-backed config
	config.C()
	b := []byte(config.Val.Gcp.Secret)

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return err
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit this URL, authorize, and paste the code here:\n%v\n", authURL)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter authorization code: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	code = strings.TrimSpace(code)

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	f, err := os.Create(tokenOutputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(tok); err != nil {
		return err
	}

	fmt.Printf("Token saved to %s\n", tokenOutputPath)
	return nil
}

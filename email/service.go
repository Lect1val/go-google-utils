package email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/Lect1val/go-google-utils/config"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type emailService struct {
	srv    *gmail.Service
	logger *zap.Logger
}

const tokenFilePath = "token.json"

// persistingTokenSource wraps an oauth2.TokenSource, ensuring any refreshed
// tokens are written to disk so subsequent runs reuse the latest token.
type persistingTokenSource struct {
	inner           oauth2.TokenSource
	path            string
	logger          *zap.Logger
	refreshFallback string
}

func (p *persistingTokenSource) Token() (*oauth2.Token, error) {
	tok, err := p.inner.Token()
	if err != nil {
		return nil, err
	}

	// Google often omits refresh_token on refresh; persist the original one.
	if tok.RefreshToken == "" && p.refreshFallback != "" {
		tok.RefreshToken = p.refreshFallback
	}

	// Save token to disk atomically.
	b, err := json.Marshal(tok)
	if err == nil {
		// Best-effort write; log on error but don't fail the request path.
		if writeErr := os.WriteFile(p.path, b, 0o600); writeErr != nil {
			p.logger.Warn("failed to persist refreshed OAuth token", zap.Error(writeErr))
		}
	}

	return tok, nil
}

func loadToken(logger *zap.Logger) (*oauth2.Token, string, error) {
	// Prefer file if present so we can reuse refreshed tokens across restarts.
	// if _, err := os.Stat(tokenFilePath); err == nil {
	// 	b, readErr := os.ReadFile(tokenFilePath)
	// 	if readErr == nil {
	// 		var t oauth2.Token
	// 		if err := json.Unmarshal(b, &t); err == nil {
	// 			return &t, t.RefreshToken, nil
	// 		} else {
	// 			logger.Warn("failed to parse token.json; falling back to env", zap.Error(err))
	// 		}
	// 	}
	// }

	// Fall back to env variable if provided.
	if config.Val.Gcp.Token != "" {
		var t oauth2.Token
		if err := json.Unmarshal([]byte(config.Val.Gcp.Token), &t); err != nil {
			return nil, "", err
		}
		return &t, t.RefreshToken, nil
	}

	return nil, "", os.ErrNotExist
}

func NewEmailService(logger *zap.Logger) EmailService {
	ctx := context.Background()

	// Ensure config is loaded
	config.C()

	configOauth, err := google.ConfigFromJSON([]byte(config.Val.Gcp.Secret), gmail.GmailSendScope)
	if err != nil {
		logger.Error(err.Error())
	}

	// Load token from file (preferred) or env, and create a persisting TokenSource
	// so refreshed tokens are saved to disk.
	token, refreshFallback, err := loadToken(logger)
	if err != nil {
		logger.Error("failed to load OAuth token; ensure initial authorization was completed", zap.Error(err))
	}

	ts := configOauth.TokenSource(ctx, token)
	pts := &persistingTokenSource{inner: ts, path: tokenFilePath, logger: logger, refreshFallback: refreshFallback}
	client := oauth2.NewClient(ctx, pts)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logger.Error(err.Error())
	}
	return &emailService{
		srv:    service,
		logger: logger,
	}
}

func (s *emailService) SendIndividualEmail(email, subject, contentType string, message string) error {
	rawMessage := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: " + contentType + "; charset=UTF-8\r\n" +
		"\r\n" +
		message)

	raw := base64.URLEncoding.EncodeToString([]byte(rawMessage))
	raw = strings.ReplaceAll(raw, "+", "-")
	raw = strings.ReplaceAll(raw, "/", "_")
	msg := &gmail.Message{
		Raw: raw,
	}
	_, err := s.srv.Users.Messages.Send("me", msg).Do()
	return err
}

package email

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func NewEmailService(logger *zap.Logger) EmailService {
	ctx := context.Background()

	// Ensure config is loaded
	config.C()

	configOauth, err := google.ConfigFromJSON([]byte(config.Val.Gcp.Secret), gmail.GmailSendScope)
	if err != nil {
		logger.Error(err.Error())
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(config.Val.Gcp.Token), &token)
	if err != nil {
		logger.Error(err.Error())
	}

	client := configOauth.Client(ctx, &token)

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

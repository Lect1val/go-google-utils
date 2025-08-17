package email

type EmailService interface {
	SendIndividualEmail(email, subject, contentType string, message string) error
}

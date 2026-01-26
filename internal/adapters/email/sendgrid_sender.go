package email

import "github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"

type sendgridSender struct {
	apiKey    string
	fromEmail string
}

// SendgridConfig holds Sendgrid configuration
type SendgridConfig struct {
	APIKey    string
	FromEmail string
}

// NewSendgridSender creates a new Sendgrid email sender
func NewSendgridSender(config SendgridConfig) auth.EmailSender {
	return &sendgridSender{
		apiKey:    config.APIKey,
		fromEmail: config.FromEmail,
	}
}

func (s *sendgridSender) Send(email auth.Email) error {
	// Implementation for sending email via SendGrid would go here
	// For now, this is a stub

	return nil
}

func (s *sendgridSender) SendWelcomeEmail(to string) error {
	email := auth.Email{
		To:      to,
		Subject: "Welcome to Tiny Store",
		Text:    "Thank you for joining Tiny Store!",
		HTML:    "<h1>Welcome to Tiny Store!</h1><p>Thank you for joining Tiny Store!</p>",
	}
	return s.Send(email)
}

func (s *sendgridSender) SendPasswordResetEmail(to, resetURL string) error {
	email := auth.Email{
		To:      to,
		Subject: "Password Reset Request",
		Text:    "Click the link to reset your password: " + resetURL,
		HTML:    "<p>Click the link to reset your password: <a href=\"" + resetURL + "\">Reset Password</a></p>",
	}
	return s.Send(email)
}

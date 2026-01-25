package integrations

import "github.com/RubenRodrigo/go-tiny-store/internal/application/ports"

type EmailService struct {
	sender ports.EmailSender
}

func NewEmailService(sender ports.EmailSender) *EmailService {
	return &EmailService{
		sender: sender,
	}
}

func (e *EmailService) SendWelcomeEmail(to string) error {
	email := ports.Email{
		To:      to,
		Subject: "Welcome to Tiny Store",
		Text:    "Thank you for joining Tiny Store!",
		HTML:    "<h1>Welcome to Tiny Store!</h1><p>Thank you for joining Tiny Store!</p>",
	}
	return e.sender.Send(email)
}

func (e *EmailService) SendPasswordResetEmail(to, resetUrl string) error {
	email := ports.Email{
		To:      to,
		Subject: "Password Reset Request",
		Text:    "Click the link to reset your password: " + resetUrl,
		HTML:    "<p>Click the link to reset your password: <a href=\"" + resetUrl + "\">Reset Password</a></p>",
	}
	return e.sender.Send(email)
}

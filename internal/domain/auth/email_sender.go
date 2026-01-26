package auth

// Email represents an email message (domain value object)
type Email struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(email Email) error
	SendWelcomeEmail(to string) error
	SendPasswordResetEmail(to, resetURL string) error
}

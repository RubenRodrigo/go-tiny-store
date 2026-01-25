package email

import "github.com/RubenRodrigo/go-tiny-store/internal/application/ports"

type SendgridManager struct {
	apiKey    string
	fromEmail string
}

type SendgridConfig struct {
	APIKey    string
	FromEmail string
}

func NewSendgridManager(config SendgridConfig) *SendgridManager {
	return &SendgridManager{
		apiKey:    config.APIKey,
		fromEmail: config.FromEmail,
	}
}

func (s *SendgridManager) Send(email ports.Email) error {
	// Implementation for sending email via SendGrid would go here
	return nil
}

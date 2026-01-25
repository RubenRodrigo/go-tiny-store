package ports

type Email struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

type EmailSender interface {
	Send(email Email) error
}

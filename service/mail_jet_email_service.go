package service

import (
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"log/slog"
)

// MailJetEmailService is an implementation of EmailService
type MailJetEmailService struct {
	client *mailjet.Client
}

// NewMailJetEmailService initializes and returns a new MailJetEmailService.
func NewMailJetEmailService() *MailJetEmailService {
	client := mailjet.NewMailjetClient("bfdda4623b8845dd83d1cf0a7b874713", "d60cddf038658d23fa21d25e481c4dcb")
	return &MailJetEmailService{client: client}
}

// SendEmail sends an email to the given email address with the given subject and body.
func (m *MailJetEmailService) SendEmail(to, subject, text string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "kontest@ayushsinghal.tech",
				Name:  "Ayush Singhal",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
				},
			},
			Subject:  subject,
			TextPart: text,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.client.SendMailV31(&messages)

	if err != nil {
		slog.Error("Error sending email: ", err)
	}

	return err
}

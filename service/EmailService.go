package service

import "fmt"

type EmailService interface {
	SendEmail(email, subject, body string) error
}

// MailJetEmailService is a mock implementation of EmailService
type MailJetEmailService struct{}

// SendEmail mocks the process of sending an email
func (m *MailJetEmailService) SendEmail(email, subject, body string) error {
	// In real implementation, integrate with MailJet API to send an email
	fmt.Printf("Sending email to: %s\nSubject: %s\nBody: %s\n", email, subject, body)
	return nil
}

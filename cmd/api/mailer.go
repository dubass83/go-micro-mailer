package api

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		contentType mail.ContentType,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type MailTrapSender struct {
	name             string
	fromEmailAdress  string
	mailtrapLogin    string
	mailtrapPassword string
	mailTrapSMTPHost string
	mailTrapSMTPPort int
	mailTrapSMTPAuth mail.SMTPAuthType
}

func NewMailtrapSender(name, email, login, pass, smtpHost string, smtpPort int, smtpAuth mail.SMTPAuthType) EmailSender {
	return &MailTrapSender{
		name:             name,
		fromEmailAdress:  email,
		mailtrapLogin:    login,
		mailtrapPassword: pass,
		mailTrapSMTPHost: smtpHost,
		mailTrapSMTPPort: smtpPort,
		mailTrapSMTPAuth: smtpAuth,
	}
}

func (sender *MailTrapSender) SendEmail(
	subject string,
	content string,
	contentType mail.ContentType,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	m := mail.NewMsg()
	if err := m.FromFormat(sender.name, sender.fromEmailAdress); err != nil {
		return fmt.Errorf("failed to set from address: %s", err)
	}
	if err := m.To(to...); err != nil {
		return fmt.Errorf("failed to set To address: %s", err)
	}
	if err := m.Cc(cc...); err != nil {
		return fmt.Errorf("failed to set CC address: %s", err)
	}
	if err := m.Bcc(bcc...); err != nil {
		return fmt.Errorf("failed to set BCC address: %s", err)
	}
	m.Subject(subject)
	// mail.TypeTextHTML
	m.SetBodyString(contentType, content)

	for _, file := range attachFiles {
		m.AttachFile(file)
	}

	c, err := mail.NewClient(
		sender.mailTrapSMTPHost,
		mail.WithPort(sender.mailTrapSMTPPort),
		mail.WithSMTPAuth(sender.mailTrapSMTPAuth),
		mail.WithUsername(sender.mailtrapLogin),
		mail.WithPassword(sender.mailtrapPassword),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}

	return c.DialAndSend(m)
}

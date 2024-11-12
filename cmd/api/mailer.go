package api

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dubass83/go-micro-mailer/util"
	"github.com/vanng822/go-premailer/premailer"
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

func NewMailtrapSender(name, email, login, pass string) EmailSender {
	return &MailTrapSender{
		name:             name,
		fromEmailAdress:  email,
		mailtrapLogin:    login,
		mailtrapPassword: pass,
		mailTrapSMTPHost: "sandbox.smtp.mailtrap.io",
		mailTrapSMTPPort: 2525,
		mailTrapSMTPAuth: mail.SMTPAuthPlain,
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

func buildHTMLMessage(conf util.Config, message map[string]any) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s", conf.TemplateDir, conf.TemplateHTML)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", fmt.Errorf("failed to create template from %s: %s", templateToRender, err)
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", message); err != nil {
		return "", fmt.Errorf("failed execute template with message %v: %s", message, err)
	}

	formattedMessage, err := inlineCSS(tpl.String())
	if err != nil {
		return "", fmt.Errorf("failed generate inline CSS message from template: %s", err)
	}
	return formattedMessage, nil
}

func inlineCSS(fm string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(fm, &options)
	if err != nil {
		return "", fmt.Errorf("failed create premailer from string %s: %s", fm, err)
	}

	html, err := prem.Transform()
	if err != nil {
		return "", fmt.Errorf("failed transform premailer to string: %s", err)
	}
	return html, nil
}

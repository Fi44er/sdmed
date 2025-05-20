package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

// Mailer структура для управления отправкой писем
type Mailer struct {
	SMTPClient *email.Pool
	Template   *template.Template
}

// NewMailer инициализирует новый Mailer с пулом соединений и шаблоном
func NewMailer(smtpHost, smtpPort, username, password, templatePath string, poolSize int) (*Mailer, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email template: %w", err)
	}

	auth := smtp.PlainAuth("", username, password, smtpHost)
	pool, err := email.NewPool(fmt.Sprintf("%s:%s", smtpHost, smtpPort), poolSize, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create email pool: %w", err)
	}

	return &Mailer{
		SMTPClient: pool,
		Template:   tmpl,
	}, nil
}

// SendMail отправляет письмо с кодом подтверждения
func (m *Mailer) SendMail(from, to, subject string, templateData interface{}) error {
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject

	var body bytes.Buffer
	if err := m.Template.Execute(&body, templateData); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	e.HTML = body.Bytes()

	if err := m.SMTPClient.Send(e, 10*time.Second); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// SendMailAsync отправляет письмо асинхронно
func (m *Mailer) SendMailAsync(from, to, subject string, templateData interface{}) {
	go func() {
		if err := m.SendMail(from, to, subject, templateData); err != nil {
			log.Printf("Error sending email to %s: %v", to, err)
		}
	}()
}

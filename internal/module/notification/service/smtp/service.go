package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/Fi44er/sdmedik/backend/internal/module/notification/service"
	"github.com/jordan-wright/email"
)

type SMTPNotifier struct {
	SMTPClient *email.Pool
	Template   *template.Template
	from       string
}

func NewSMTPNotifier(smtpHost, smtpPort, username, password, templatePath string, poolSize int) (*SMTPNotifier, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email template: %w", err)
	}

	auth := smtp.PlainAuth("", username, password, smtpHost)
	pool, err := email.NewPool(fmt.Sprintf("%s:%s", smtpHost, smtpPort), poolSize, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create email pool: %w", err)
	}

	return &SMTPNotifier{
		SMTPClient: pool,
		Template:   tmpl,
		from:       username,
	}, nil
}

func (n *SMTPNotifier) Send(msg *service.Message) error {
	e := email.NewEmail()
	e.From = n.from
	e.To = []string{msg.Recipient}
	e.Subject = msg.Subject

	var body bytes.Buffer
	if err := n.Template.Execute(&body, msg.Data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	e.HTML = body.Bytes()

	if err := n.SMTPClient.Send(e, 10*time.Second); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

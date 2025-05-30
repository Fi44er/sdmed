package smtp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/Fi44er/sdmedik/backend/internal/module/notification/service"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jordan-wright/email"
)

type SMTPNotifier struct {
	SMTPClient *email.Pool
	from       string
}

func NewSMTPNotifier(smtpHost, smtpPort, username, password string, poolSize int) (*SMTPNotifier, error) {
	auth := smtp.PlainAuth("", username, password, smtpHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,     // ⚠️ ОПАСНО: отключает проверку сертификата
		ServerName:         smtpHost, // Укажите ваш SMTP-хост
	}
	pool, err := email.NewPool(fmt.Sprintf("%s:%s", smtpHost, smtpPort), poolSize, auth, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create email pool: %w", err)
	}

	return &SMTPNotifier{
		SMTPClient: pool,
		from:       username,
	}, nil
}

func (n *SMTPNotifier) Send(msg *service.Message) {
	go func() {
		if err := n.send(msg); err != nil {
			log.Infof("Error sending email to: %v", err)
		}
	}()

}

func (n *SMTPNotifier) send(msg *service.Message) error {
	e := email.NewEmail()
	e.From = n.from
	e.To = []string{msg.Recipient}
	e.Subject = msg.Subject

	if msg.TemplatePath != "" {
		tmpl, err := template.ParseFiles(msg.TemplatePath)
		if err != nil {
			return fmt.Errorf("failed to parse email template: %w", err)
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, msg.Data); err != nil {
			return fmt.Errorf("failed to execute email template: %w", err)
		}
		e.HTML = body.Bytes()
	} else if msg.Content != "" {
		e.HTML = []byte(msg.Content)
	} else {
		return fmt.Errorf("either TemplatePath or Content must be provided")
	}

	if err := n.SMTPClient.Send(e, 10*time.Second); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

package pkg

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

type Mailer struct {
	client   *mail.Client
	fromAddr string
	fromName string
}

func NewMailer(host string, port int, user string, pass string, fromName string) (*Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(user),
		mail.WithPassword(pass),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	return &Mailer{
		client:   client,
		fromAddr: user,
		fromName: fromName,
	}, nil
}

func (m *Mailer) Send(to string, subject string, bodyHTML string) error {
	msg := mail.NewMsg()

	if err := msg.FromFormat(m.fromName, m.fromAddr); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}

	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, bodyHTML)

	if err := m.client.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

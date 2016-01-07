package mailer

import (
	"net/smtp"
	"strings"
)

type Mailer interface {
	Send(string, []string, string, string) error
}

func New(addr string, auth smtp.Auth) Mailer {
	return &defaultMailer{addr, auth}
}

type defaultMailer struct {
	addr string
	auth smtp.Auth
}

func (m *defaultMailer) Send(from string, to []string, subject, message string) error {

	parts := []string{
		"To: " + strings.Join(to, ","),
		"From: " + from,
		"Subject: " + subject,
		"",
		message,
		"",
	}
	return smtp.SendMail(
		m.addr,
		m.auth,
		from,
		to,
		[]byte(strings.Join(parts, "\r\n")),
	)
}

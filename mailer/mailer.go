package mailer

import (
	"bytes"
	"net/smtp"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/danjac/podbaby/config"
)

// Mailer handles SMTP calls
type Mailer interface {
	// Send sends a plain SMTP email
	Send(string, []string, string, string) error
	// SendFromTemplate renders message to template and sends plain SMTP email
	SendFromTemplate(string, []string, string, string, interface{}) error
}

// New returns a new Mailer instance
func New(cfg *config.Config) (Mailer, error) {

	auth := smtp.PlainAuth(
		cfg.Mail.ID,
		cfg.Mail.User,
		cfg.Mail.Password,
		cfg.Mail.Host,
	)

	templates, err := template.ParseGlob(filepath.Join(cfg.Mail.TemplateDir, "*.tmpl"))
	if err != nil {
		return nil, err
	}
	return &defaultMailer{cfg.Mail.Addr, auth, templates}, nil
}

type defaultMailer struct {
	addr      string
	auth      smtp.Auth
	templates *template.Template
}

func (m *defaultMailer) SendFromTemplate(from string, to []string, subject string, template string, data interface{}) error {
	var buf bytes.Buffer
	err := m.templates.ExecuteTemplate(&buf, template, data)
	if err != nil {
		return err
	}
	return m.Send(from, to, subject, buf.String())
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

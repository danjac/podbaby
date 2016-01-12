package mailer

import (
	"bytes"
	"net/smtp"
	"path/filepath"
	"strings"
	"text/template"
)

type Mailer interface {
	Send(string, []string, string, string) error
	SendFromTemplate(string, []string, string, string, interface{}) error
}

func New(addr string, auth smtp.Auth, templateDir string) (Mailer, error) {
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.tmpl"))
	if err != nil {
		return nil, err
	}
	return &defaultMailer{addr, auth, templates}, nil
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

package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/open-pm/open-pm/server/internal/config"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*.html
var templateFS embed.FS

// Mailer sends emails via SMTP.
type Mailer struct {
	host     string
	port     int
	user     string
	pass     string
	fromName string
	fromAddr string
	tmpl     *template.Template
}

// New creates a Mailer from config.
func New(cfg config.SMTPConfig) (*Mailer, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("email: failed to parse templates: %w", err)
	}

	return &Mailer{
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		pass:     cfg.Pass,
		fromName: cfg.FromName,
		fromAddr: cfg.FromAddr,
		tmpl:     tmpl,
	}, nil
}

// Send renders a template and sends it via SMTP.
func (m *Mailer) Send(to, subject, templateName string, data interface{}) error {
	var body bytes.Buffer
	if err := m.tmpl.ExecuteTemplate(&body, templateName, data); err != nil {
		return fmt.Errorf("email: failed to render template %q: %w", templateName, err)
	}

	from := fmt.Sprintf("%s <%s>", m.fromName, m.fromAddr)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, body.String())

	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}

	if err := smtp.SendMail(addr, auth, m.fromAddr, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("email: failed to send to %s: %w", to, err)
	}

	log.Debug().Str("to", to).Str("subject", subject).Msg("email sent")
	return nil
}

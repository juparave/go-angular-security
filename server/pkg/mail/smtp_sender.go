package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPSender implements EmailSender using SMTP
type SMTPSender struct {
	config *Config
}

// NewSMTPSender creates a new SMTP email sender
func NewSMTPSender(cfg *Config) *SMTPSender {
	return &SMTPSender{config: cfg}
}

// Send sends a plain text email
func (s *SMTPSender) Send(to, subject, body string) error {
	return s.sendEmail(to, subject, body, "text/plain")
}

// SendHTML sends an HTML email
func (s *SMTPSender) SendHTML(to, subject, htmlBody string) error {
	return s.sendEmail(to, subject, htmlBody, "text/html")
}

// sendEmail sends an email via SMTP
func (s *SMTPSender) sendEmail(to, subject, body, contentType string) error {
	if s.config.Host == "" {
		return &EmailError{Message: "SMTP host not configured"}
	}

	// Build the email message
	from := s.config.FromAddress
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress)
	}

	msg := bytes.NewBuffer(nil)
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString(fmt.Sprintf("MIME-Version: 1.0\r\n"))
	msg.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\r\n", contentType))
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	var auth smtp.Auth
	if s.config.Account != "" && s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.Account, s.config.Password, s.config.Host)
	}

	// Determine if we need TLS
	var client *smtp.Client
	var err error

	if s.config.Port == "465" {
		// Implicit TLS (SMTPS)
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.config.Host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return &EmailError{Message: "failed to connect via TLS", Cause: err}
		}
		client, err = smtp.NewClient(conn, s.config.Host)
		if err != nil {
			return &EmailError{Message: "failed to create SMTP client", Cause: err}
		}
	} else {
		// Standard SMTP (potentially with STARTTLS)
		client, err = smtp.Dial(addr)
		if err != nil {
			return &EmailError{Message: "failed to connect to SMTP server", Cause: err}
		}
	}

	defer client.Close()

	// Send HELO/EHLO
	if err := client.Hello("localhost"); err != nil {
		return &EmailError{Message: "HELO/EHLO failed", Cause: err}
	}

	// STARTTLS for port 587
	if s.config.Port == "587" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.config.Host,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return &EmailError{Message: "STARTTLS failed", Cause: err}
		}
	}

	// Authenticate if credentials provided
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return &EmailError{Message: "SMTP authentication failed", Cause: err}
		}
	}

	// Set sender
	if err := client.Mail(s.config.FromAddress); err != nil {
		return &EmailError{Message: "failed to set sender", Cause: err}
	}

	// Set recipients
	for _, recipient := range strings.Split(to, ",") {
		recipient = strings.TrimSpace(recipient)
		if err := client.Rcpt(recipient); err != nil {
			return &EmailError{Message: "failed to set recipient", Cause: err}
		}
	}

	// Send data
	w, err := client.Data()
	if err != nil {
		return &EmailError{Message: "failed to prepare email data", Cause: err}
	}

	_, err = w.Write(msg.Bytes())
	if err != nil {
		return &EmailError{Message: "failed to write email data", Cause: err}
	}

	if err := w.Close(); err != nil {
		return &EmailError{Message: "failed to close email data", Cause: err}
	}

	return client.Quit()
}

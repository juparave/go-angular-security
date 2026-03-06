package handlers

import (
	"server/internal/config"
	"server/pkg/mail"
)

var app *config.AppConfig
var emailService *mail.Service

func SetRepo(a *config.AppConfig) {
	app = a
}

// SetEmailService sets the email service for sending emails
func SetEmailService(svc *mail.Service) {
	emailService = svc
}

// GetEmailService returns the email service
func GetEmailService() *mail.Service {
	return emailService
}

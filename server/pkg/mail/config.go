package mail

// Config holds the email service configuration
type Config struct {
	Host          string // SMTP server host
	Port          string // SMTP server port
	Account       string // SMTP authentication account/username
	Password      string // SMTP authentication password
	FromAddress   string // Default "from" email address
	FromName      string // Default "from" display name
	DevMode       bool   // When true, redirect all emails to TestRecipient
	TestRecipient string // Email address to receive all emails in dev mode
}

// NewConfig creates a new email configuration
func NewConfig(host, port, account, password, fromAddress, fromName string) *Config {
	return &Config{
		Host:        host,
		Port:        port,
		Account:     account,
		Password:    password,
		FromAddress: fromAddress,
		FromName:    fromName,
	}
}

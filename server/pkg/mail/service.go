package mail

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// EmailMessage represents an email to be sent
type EmailMessage struct {
	To       string
	Subject  string
	Template string
	Data     map[string]interface{}
}

// Service provides async email sending capabilities
type Service struct {
	config    *Config
	sender    EmailSender
	templates *TemplateManager

	// Async processing
	emailChan chan EmailMessage
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc

	// Metrics
	sentCount atomic.Int64
	failCount atomic.Int64
	queueSize int
}

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(to, subject, body string) error
	SendHTML(to, subject, htmlBody string) error
}

// NewService creates a new email service
func NewService(cfg *Config, queueSize int) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Service{
		config:    cfg,
		sender:    NewSMTPSender(cfg),
		templates: NewTemplateManager(),
		emailChan: make(chan EmailMessage, queueSize),
		ctx:       ctx,
		cancel:    cancel,
		queueSize: queueSize,
	}

	// Start worker goroutines
	s.startWorkers(3)

	return s
}

// startWorkers starts the email worker goroutines
func (s *Service) startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// worker processes emails from the queue
func (s *Service) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case email, ok := <-s.emailChan:
			if !ok {
				return
			}
			if err := s.sendSync(email); err != nil {
				s.failCount.Add(1)
				log.Printf("Failed to send email to %s: %v", email.To, err)
			} else {
				s.sentCount.Add(1)
				log.Printf("Email sent successfully to %s", email.To)
			}
		}
	}
}

// Send queues an email to be sent asynchronously
func (s *Service) Send(to, subject, templateName string, data map[string]interface{}) error {
	email := EmailMessage{
		To:       to,
		Subject:  subject,
		Template: templateName,
		Data:     data,
	}

	select {
	case s.emailChan <- email:
		return nil
	default:
		return ErrQueueFull
	}
}

// sendSync renders and sends an email synchronously with up to 3 retry attempts.
func (s *Service) sendSync(email EmailMessage) error {
	to := email.To
	// Dev mode: redirect all emails to the test recipient
	if s.config.DevMode && s.config.TestRecipient != "" {
		log.Printf("Dev mode: redirecting email for %s to %s", email.To, s.config.TestRecipient)
		to = s.config.TestRecipient
	}

	// Render template
	body, err := s.templates.Render(email.Template, email.Data)
	if err != nil {
		return err
	}

	// Retry up to 3 times with exponential backoff (1s, 2s, 4s)
	var lastErr error
	backoff := time.Second
	for attempt := 1; attempt <= 3; attempt++ {
		lastErr = s.sender.SendHTML(to, email.Subject, body)
		if lastErr == nil {
			return nil
		}
		log.Printf("Email send attempt %d/3 failed for %s: %v", attempt, to, lastErr)
		if attempt < 3 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return lastErr
}

// SendRaw sends a plain text email synchronously
func (s *Service) SendRaw(to, subject, body string) error {
	return s.sender.Send(to, subject, body)
}

// SendHTML sends an HTML email synchronously without template
func (s *Service) SendHTML(to, subject, htmlBody string) error {
	return s.sender.SendHTML(to, subject, htmlBody)
}

// Close stops the email service and waits for pending emails
func (s *Service) Close() error {
	s.cancel()
	close(s.emailChan)
	s.wg.Wait()
	return nil
}

// Metrics returns the current email metrics
func (s *Service) Metrics() Metrics {
	return Metrics{
		Sent:     s.sentCount.Load(),
		Failed:   s.failCount.Load(),
		Queued:   len(s.emailChan),
		Capacity: s.queueSize,
	}
}

// Metrics holds email service metrics
type Metrics struct {
	Sent     int64
	Failed   int64
	Queued   int
	Capacity int
}

// Errors
var (
	ErrQueueFull = &EmailError{Message: "email queue is full"}
)

// EmailError represents an email-related error
type EmailError struct {
	Message string
	Cause   error
}

func (e *EmailError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *EmailError) Unwrap() error {
	return e.Cause
}

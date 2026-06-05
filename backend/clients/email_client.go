package clients

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

// TicketInfo represents the ticket data needed by email templates.
// It is intentionally decoupled from domain.Ticket to avoid
// introducing a dependency on the domain layer.
type TicketInfo struct {
	TicketID   uint
	EventTitle string
	EventDate  string
	Location   string
	Price      float64
}

// EmailClient defines the interface for sending transactional emails.
// Services depend on this interface so they remain testable with mocks.
type EmailClient interface {
	SendPurchaseConfirmation(to string, ticket TicketInfo) error
	SendCancellationNotice(to string, ticket TicketInfo) error
	SendTransferNotice(from, to string, ticket TicketInfo) error
}

// NewEmailClient returns the appropriate EmailClient implementation
// based on the EMAIL_PROVIDER environment variable.
//
// Supported providers:
//
//	"smtp" – sends via an SMTP server (configure SMTP_HOST, SMTP_PORT, etc.)
//	""      – logs to stdout (development default, no external dependency)
func NewEmailClient() EmailClient {
	provider := os.Getenv("EMAIL_PROVIDER")
	switch provider {
	case "smtp":
		return &smtpEmailClient{
			host:     os.Getenv("SMTP_HOST"),
			port:     os.Getenv("SMTP_PORT"),
			username: os.Getenv("SMTP_USERNAME"),
			password: os.Getenv("SMTP_PASSWORD"),
			from:     os.Getenv("EMAIL_FROM"),
		}
	default:
		return &logEmailClient{}
	}
}

// ---------------------------------------------------------------------------
// Log client (development / fallback)
// ---------------------------------------------------------------------------

type logEmailClient struct{}

func (c *logEmailClient) SendPurchaseConfirmation(to string, ticket TicketInfo) error {
	log.Printf("[EMAIL] Purchase confirmation → %s | event=%s ticket=#%d", to, ticket.EventTitle, ticket.TicketID)
	return nil
}

func (c *logEmailClient) SendCancellationNotice(to string, ticket TicketInfo) error {
	log.Printf("[EMAIL] Cancellation notice → %s | event=%s ticket=#%d", to, ticket.EventTitle, ticket.TicketID)
	return nil
}

func (c *logEmailClient) SendTransferNotice(from, to string, ticket TicketInfo) error {
	log.Printf("[EMAIL] Transfer notice → %s (from %s) | event=%s ticket=#%d", to, from, ticket.EventTitle, ticket.TicketID)
	return nil
}

// ---------------------------------------------------------------------------
// SMTP client (production)
// ---------------------------------------------------------------------------

type smtpEmailClient struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func (c *smtpEmailClient) SendPurchaseConfirmation(to string, ticket TicketInfo) error {
	subject := "Purchase Confirmed – Ceibo Tickets"
	body := fmt.Sprintf(`Hi,

Your purchase for %s has been confirmed.

  Event:     %s
  Date:      %s
  Location:  %s
  Ticket #:  %d
  Price:     $%.2f

Enjoy the event!
Ceibo Tickets`, ticket.EventTitle, ticket.EventTitle, ticket.EventDate, ticket.Location, ticket.TicketID, ticket.Price)

	return c.send(to, subject, body)
}

func (c *smtpEmailClient) SendCancellationNotice(to string, ticket TicketInfo) error {
	subject := "Ticket Cancelled – Ceibo Tickets"
	body := fmt.Sprintf(`Hi,

Your ticket for %s has been cancelled.

  Event:     %s
  Date:      %s
  Ticket #:  %d

We hope to see you at future events.
Ceibo Tickets`, ticket.EventTitle, ticket.EventTitle, ticket.EventDate, ticket.TicketID)

	return c.send(to, subject, body)
}

func (c *smtpEmailClient) SendTransferNotice(from, to string, ticket TicketInfo) error {
	subject := "Ticket Transferred – Ceibo Tickets"
	body := fmt.Sprintf(`Hi,

A ticket for %s has been transferred to you from %s.

  Event:     %s
  Date:      %s
  Ticket #:  %d

Enjoy the event!
Ceibo Tickets`, ticket.EventTitle, from, ticket.EventTitle, ticket.EventDate, ticket.TicketID)

	return c.send(to, subject, body)
}

func (c *smtpEmailClient) send(to, subject, body string) error {
	const sendTimeout = 30 * time.Second

	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	auth := smtp.PlainAuth("", c.username, c.password, c.host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", c.from, to, subject, body)

	errCh := make(chan error, 1)
	go func() {
		errCh <- smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg))
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("email client: send mail: %w", err)
		}
		return nil
	case <-time.After(sendTimeout):
		return errors.New("email client: send mail timed out")
	}
}

// Ensure compile-time interface satisfaction.
var _ EmailClient = (*logEmailClient)(nil)
var _ EmailClient = (*smtpEmailClient)(nil)

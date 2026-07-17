package clients

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailClientDefault(t *testing.T) {
	os.Unsetenv("EMAIL_PROVIDER")
	client := NewEmailClient()
	_, ok := client.(*logEmailClient)
	assert.True(t, ok, "default provider should return logEmailClient")
}

func TestNewEmailClientLog(t *testing.T) {
	os.Setenv("EMAIL_PROVIDER", "log")
	defer os.Unsetenv("EMAIL_PROVIDER")
	client := NewEmailClient()
	_, ok := client.(*logEmailClient)
	assert.True(t, ok)
}

func TestNewEmailClientSMTP(t *testing.T) {
	os.Setenv("EMAIL_PROVIDER", "smtp")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "user")
	os.Setenv("SMTP_PASSWORD", "pass")
	os.Setenv("EMAIL_FROM", "noreply@example.com")
	defer func() {
		os.Unsetenv("EMAIL_PROVIDER")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_USERNAME")
		os.Unsetenv("SMTP_PASSWORD")
		os.Unsetenv("EMAIL_FROM")
	}()

	client := NewEmailClient()
	smtpClient, ok := client.(*smtpEmailClient)
	assert.True(t, ok)
	assert.Equal(t, "smtp.example.com", smtpClient.host)
	assert.Equal(t, "587", smtpClient.port)
	assert.Equal(t, "user", smtpClient.username)
	assert.Equal(t, "pass", smtpClient.password)
	assert.Equal(t, "noreply@example.com", smtpClient.from)
}

func TestLogEmailClient_SendPurchaseConfirmation(t *testing.T) {
	c := &logEmailClient{}
	err := c.SendPurchaseConfirmation("test@test.com", TicketInfo{
		TicketID:   1,
		EventTitle: "Test Event",
	})
	assert.NoError(t, err)
}

func TestLogEmailClient_SendCancellationNotice(t *testing.T) {
	c := &logEmailClient{}
	err := c.SendCancellationNotice("test@test.com", TicketInfo{
		TicketID:   1,
		EventTitle: "Test Event",
	})
	assert.NoError(t, err)
}

func TestLogEmailClient_SendTransferNotice(t *testing.T) {
	c := &logEmailClient{}
	err := c.SendTransferNotice("from@test.com", "to@test.com", TicketInfo{
		TicketID:   1,
		EventTitle: "Test Event",
	})
	assert.NoError(t, err)
}

func TestSMTPEmailClient_SendPurchaseConfirmationInvalidHost(t *testing.T) {
	c := &smtpEmailClient{
		host:     "",
		port:     "25",
		username: "",
		password: "",
		from:     "noreply@test.com",
	}
	err := c.SendPurchaseConfirmation("test@test.com", TicketInfo{
		TicketID:   1,
		EventTitle: "Test Event",
		EventDate:  "2026-07-20",
		Location:   "Test Location",
		Price:      100.00,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email client:")
}

func TestSMTPEmailClient_SendCancellationNoticeInvalidHost(t *testing.T) {
	c := &smtpEmailClient{
		host:     "",
		port:     "25",
		username: "",
		password: "",
		from:     "noreply@test.com",
	}
	err := c.SendCancellationNotice("test@test.com", TicketInfo{
		TicketID:   1,
		EventTitle: "Test Event",
	})
	assert.Error(t, err)
}

func TestSMTPEmailClient_SendTransferNoticeInvalidHost(t *testing.T) {
	c := &smtpEmailClient{
		host:     "",
		port:     "25",
		username: "",
		password: "",
		from:     "noreply@test.com",
	}
	err := c.SendTransferNotice("from@test.com", "to@test.com", TicketInfo{
		TicketID:   1,
		EventTitle: "Test Event",
	})
	assert.Error(t, err)
}

func TestNewEmailClientUnknownProvider(t *testing.T) {
	os.Setenv("EMAIL_PROVIDER", "unknown")
	defer os.Unsetenv("EMAIL_PROVIDER")
	client := NewEmailClient()
	_, ok := client.(*logEmailClient)
	assert.True(t, ok, "unknown provider should fallback to logEmailClient")
}

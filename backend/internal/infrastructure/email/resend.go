// Package email provides a thin wrapper around the Resend API for sending
// transactional emails.
package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

// ResendSender sends emails using the Resend API.
type ResendSender struct {
	client *resend.Client
	from   string
}

// NewResendSender creates a ResendSender that authenticates with the given
// Resend API key and sends emails from the given address
// (e.g. "Acme <invites@acme.com>").
func NewResendSender(apiKey, from string) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

// Send sends a single HTML email to the given recipient.
func (s *ResendSender) Send(ctx context.Context, to, subject, htmlBody string) error {
	_, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	})
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

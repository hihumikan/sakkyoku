package notifier

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	WebhookURL string
	Username   string
}

func NewSlackNotifier(webhookURL, username string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
		Username:   username,
	}
}

func (s *SlackNotifier) Notify(message string) error {
	if s.WebhookURL == "" {
		return nil
	}

	msg := &slack.WebhookMessage{
		Username: s.Username,
		Text:     message,
	}

	err := slack.PostWebhook(s.WebhookURL, msg)
	if err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}

	return nil
}

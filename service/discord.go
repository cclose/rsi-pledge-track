package service

type DiscordService struct {
	URL          string
	WebhookID    string
	WebhookToken string
}

type ExecuteWebhookPayload struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

type SendMessageOption func()

func (s *DiscordService) SendMessage(msg string, opts ...SendMessageOption) error {
	return nil
}

// Example 
// curl -X POST \
//  -H "Content-Type: application/json" \
//  -d '{"content":"Testing Webhook", "username":"${DISCORD_USERNAME}"}' \
//  ${DISCORD_WEBHOOK}

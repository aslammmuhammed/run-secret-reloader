package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/alert"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	"github.com/slack-go/slack"
)

// SlackClient implements the alert.Client interface for Slack.
type SlackClient struct {
	webhookURL string
	logger     *logger.Logger
}

// NewAlertClient creates a new Slack alert client.
func NewAlertClient(cfg *config.Config, logger *logger.Logger) *SlackClient {
	return &SlackClient{
		webhookURL: cfg.Alert.Slack.WebhookURL,
		logger:     logger,
	}
}

// SendMessage sends a formatted message to a Slack webhook.
func (s *SlackClient) SendMessage(ctx context.Context, msg alert.Message) error {
	blocks := s.buildMessageBlocks(msg)

	webhookMsg := slack.WebhookMessage{
		Blocks: &slack.Blocks{
			BlockSet: blocks,
		},
	}

	err := slack.PostWebhook(s.webhookURL, &webhookMsg)
	if err != nil {
		s.logger.Error(ctx, fmt.Sprintf("failed to send slack message: %v", err))
		return err
	}

	s.logger.Info(ctx, "successfully sent message to slack")
	return nil
}

// buildMessageBlocks constructs the Slack Block Kit structure.
func (s *SlackClient) buildMessageBlocks(msg alert.Message) []slack.Block {
	headerText := "Secret Version Update Processed"
	headerBlock := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false))

	fieldsBlock := slack.NewSectionBlock(nil, []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Secret Name:*\n`%s`", msg.SecretName), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*New Version:*\n`%s`", msg.NewVersion), false, false),
	}, nil)

	contextBlock := slack.NewContextBlock("",
		slack.NewTextBlockObject("mrkdwn", "Cloud Run Service Redeployment Status", false, false),
	)

	succeededHeaderBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf(":white_check_mark: *Successfully Redeployed Services (%d)*", len(msg.Succeeded)),
		false, false), nil, nil)
	succeededBodyBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", formatServices(msg.Succeeded), false, false), nil, nil)

	failedHeaderBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf(":x: *Failed to Redeploy Services (%d)*", len(msg.Failed)),
		false, false), nil, nil)
	failedBodyBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", formatServices(msg.Failed), false, false), nil, nil)

	skippedHeaderBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf(":arrow_right: *Skipped Services (Already Up-to-date) (%d)*", len(msg.Skipped)),
		false, false), nil, nil)
	skippedBodyBlock := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", formatServices(msg.Skipped), false, false), nil, nil)

	traceBlock := slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("Trace ID: `%s`", msg.TraceID), false, false))

	return []slack.Block{
		headerBlock,
		fieldsBlock,
		contextBlock,
		slack.NewDividerBlock(),
		succeededHeaderBlock,
		succeededBodyBlock,
		failedHeaderBlock,
		failedBodyBlock,
		skippedHeaderBlock,
		skippedBodyBlock,
		slack.NewDividerBlock(),
		traceBlock,
	}
}

func formatServices(services []string) string {
	if len(services) == 0 {
		return "_None_"
	}
	return "• `" + strings.Join(services, "`\n• `") + "`"
}

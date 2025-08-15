package slack

import (
	"github.com/aslammmuhammed/run-secret-reloader/config"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/alert"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
)

// RegisterProvider registers the Slack provider with the alert factory.
func RegisterProvider(factory *alert.Factory, cfg *config.Config, logger *logger.Logger) {
	factory.Register(alert.AlertTypeSlack, func() (alert.Client, error) {
		return NewAlertClient(cfg, logger), nil
	})
}

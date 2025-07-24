package usecase

import (
	"context"
	"strconv"

	appError "github.com/aslammmuhammed/run-secret-reloader/internal/errors"
	"github.com/aslammmuhammed/run-secret-reloader/internal/repo"
	"github.com/aslammmuhammed/run-secret-reloader/pkg/logger"
	run "google.golang.org/api/run/v1"
)

type WebhookUsecase struct {
	repo repo.CloudRunRepository
	log  logger.Logger
}

func NewWebhookUsecase(repo repo.CloudRunRepository, log logger.Logger) *WebhookUsecase {
	return &WebhookUsecase{repo: repo, log: log}
}

func (u *WebhookUsecase) GetServicesByLabel(ctx context.Context, labelKey string) ([]*run.Service, error) {

	u.log.Info(ctx, "Getting Services by labelKey: "+labelKey)

	services, err := u.repo.GetServicesByLabel(ctx, labelKey)
	if err != nil {
		u.log.Error(ctx, "Failed to get services by labelKey: "+labelKey+":"+err.Error())
		return nil, appError.NewInternalError(err.Error())
	}
	u.log.Info(ctx, "Successfully retrieved services with labelKey: "+labelKey+" | serviceCount: "+strconv.Itoa(len(services)))

	return services, nil
}

package usecase

import (
	"context"
	"strconv"
	"sync"

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

func (u *WebhookUsecase) UpdateCloudRunServices(ctx context.Context, services []*run.Service, newAnnotations map[string]string, newLabels map[string]string, newVersion string, versionAnnotationKey string) ([]*run.Service, []*run.Service, []*run.Service) {
	var wg sync.WaitGroup
	succeededChan := make(chan *run.Service, len(services))
	failedChan := make(chan *run.Service, len(services))
	skippedChan := make(chan *run.Service, len(services))

	for _, service := range services {
		// Skip if current version >= new version
		u.log.Info(ctx, "Checking if service "+service.Metadata.Name+" is up to date , oldVersion: "+service.Spec.Template.Metadata.Annotations[versionAnnotationKey]+" newVersion: "+newVersion)
		if oldVersion, ok := service.Spec.Template.Metadata.Annotations[versionAnnotationKey]; ok {
			oldVersionInt, _ := strconv.Atoi(oldVersion)
			newVersionInt, _ := strconv.Atoi(newVersion)
			if oldVersionInt >= newVersionInt {
				u.log.Info(ctx, "Skipping service "+service.Metadata.Name+" – current version "+oldVersion+">= incoming version "+newVersion)
				skippedChan <- service
				continue
			}
		}

		wg.Add(1)
		go func(s *run.Service) {
			defer wg.Done()
			err := u.repo.UpdateCloudRunAnnotationsAndLabels(ctx, s, newAnnotations, newLabels)
			if err != nil {
				u.log.Error(ctx, "Failed to update service "+s.Metadata.Name+": "+err.Error())
				failedChan <- s
			} else {
				u.log.Info(ctx, "Successfully updated service "+s.Metadata.Name)
				succeededChan <- s
			}
		}(service)
	}

	wg.Wait()
	close(succeededChan)
	close(failedChan)
	close(skippedChan)
	var succeeded []*run.Service
	for s := range succeededChan {
		succeeded = append(succeeded, s)
	}

	var failed []*run.Service
	for f := range failedChan {
		failed = append(failed, f)
	}

	var skipped []*run.Service
	for s := range skippedChan {
		skipped = append(skipped, s)
	}

	u.log.Info(ctx, "Finished updating services. Succeeded: "+strconv.Itoa(len(succeeded))+" Failed: "+strconv.Itoa(len(failed))+" Skipped: "+strconv.Itoa(len(skipped)))

	return succeeded, failed, skipped
}

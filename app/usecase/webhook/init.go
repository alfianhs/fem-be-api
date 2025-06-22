package webhook_usecase

import (
	"app/domain"
	"time"
)

type webhookAppUsecase struct {
	mongoDbRepo    domain.MongoDbRepo
	xenditRepo     domain.XenditRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDbRepo domain.MongoDbRepo
	XenditRepo  domain.XenditRepo
}

func NewWebhookAppUsecase(repoInjection RepoInjection, timeout time.Duration) domain.WebhookAppUsecase {
	return &webhookAppUsecase{
		mongoDbRepo:    repoInjection.MongoDbRepo,
		xenditRepo:     repoInjection.XenditRepo,
		contextTimeout: timeout,
	}
}

package admin_usecase

import (
	"app/domain"
	"time"
)

type adminAppUsecase struct {
	mongoDbRepo    domain.MongoDbRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDbRepo domain.MongoDbRepo
}

func NewAdminAppUsecase(repoInjection RepoInjection, timeout time.Duration) domain.AdminAppUsecase {
	return &adminAppUsecase{
		mongoDbRepo:    repoInjection.MongoDbRepo,
		contextTimeout: timeout,
	}
}

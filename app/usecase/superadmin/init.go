package superadmin_usecase

import (
	"app/domain"
	"time"
)

type superadminAppUsecase struct {
	mongoDbRepo    domain.MongoDbRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDbRepo domain.MongoDbRepo
}

func NewSuperadminAppUsecase(repoInjection RepoInjection, timeout time.Duration) domain.SuperadminAppUsecase {
	return &superadminAppUsecase{
		mongoDbRepo:    repoInjection.MongoDbRepo,
		contextTimeout: timeout,
	}
}

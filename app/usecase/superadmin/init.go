package superadmin_usecase

import (
	"app/domain"
	"time"
)

type superadminAppUsecase struct {
	mongoDbRepo    domain.MongoDbRepo
	s3Repo         domain.S3Repo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDbRepo domain.MongoDbRepo
	S3Repo      domain.S3Repo
}

func NewSuperadminAppUsecase(repoInjection RepoInjection, timeout time.Duration) domain.SuperadminAppUsecase {
	return &superadminAppUsecase{
		mongoDbRepo:    repoInjection.MongoDbRepo,
		s3Repo:         repoInjection.S3Repo,
		contextTimeout: timeout,
	}
}

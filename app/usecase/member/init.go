package member_usecase

import (
	"app/domain"
	"time"
)

type memberAppUsecase struct {
	mongoDbRepo    domain.MongoDbRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDbRepo domain.MongoDbRepo
}

func NewMemberAppUsecase(repoInjection RepoInjection, timeout time.Duration) domain.MemberAppUsecase {
	return &memberAppUsecase{
		mongoDbRepo:    repoInjection.MongoDbRepo,
		contextTimeout: timeout,
	}
}

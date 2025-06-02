package member_usecase

import (
	"app/domain"
	"time"
)

type memberAppUsecase struct {
	mongoDbRepo    domain.MongoDbRepo
	xenditRepo     domain.XenditRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDbRepo domain.MongoDbRepo
	XenditRepo  domain.XenditRepo
}

func NewMemberAppUsecase(repoInjection RepoInjection, timeout time.Duration) domain.MemberAppUsecase {
	return &memberAppUsecase{
		mongoDbRepo:    repoInjection.MongoDbRepo,
		xenditRepo:     repoInjection.XenditRepo,
		contextTimeout: timeout,
	}
}

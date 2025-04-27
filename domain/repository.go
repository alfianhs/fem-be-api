package domain

import (
	mongo_model "app/domain/model/mongo"
	"context"
)

type MongoDbRepo interface {
	// Superadmin
	FetchOneSuperadmin(ctx context.Context, options map[string]interface{}) (row *mongo_model.Superadmin, err error)

	// Admin
	FetchOneAdmin(ctx context.Context, options map[string]interface{}) (row *mongo_model.Admin, err error)

	// Member
	FetchOneMember(ctx context.Context, options map[string]interface{}) (row *mongo_model.Member, err error)
	CreateOneMember(ctx context.Context, member *mongo_model.Member) (err error)
	UpdatePartialMember(ctx context.Context, options, field map[string]interface{}) (err error)
}

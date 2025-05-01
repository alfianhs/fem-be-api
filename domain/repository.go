package domain

import (
	mongo_model "app/domain/model/mongo"
	s3_model "app/domain/model/s3"
	"context"
	"io"
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

	// Media
	FetchOneMedia(ctx context.Context, options map[string]interface{}) (row *mongo_model.Media, err error)
	CreateOneMedia(ctx context.Context, media *mongo_model.Media) (err error)
	CreateManyMedia(ctx context.Context, medias []*mongo_model.Media) (err error)
	UpdatePartialMedia(ctx context.Context, options, field map[string]interface{}) (err error)
	UpdateManyMediaPartial(ctx context.Context, options, field map[string]interface{}) (err error)
}

type S3Repo interface {
	UploadFilePublic(ctx context.Context, objectName string, body io.Reader, mimeType string) (*s3_model.UploadResponse, error)
}

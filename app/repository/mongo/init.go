package mongo_repository

import (
	"app/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDbRepo struct {
	Conn                 *mongo.Database
	superadminCollection string
	adminCollection      string
	memberCollection     string
	mediaCollection      string
}

func NewMongoDbRepo(conn *mongo.Database) domain.MongoDbRepo {
	return &mongoDbRepo{
		Conn:                 conn,
		superadminCollection: "superadmins",
		adminCollection:      "admins",
		memberCollection:     "members",
		mediaCollection:      "medias",
	}
}

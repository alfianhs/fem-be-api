package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	StageName *string            `bson:"stageName" json:"stageName"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

type PlayerFK struct {
	ID        string  `bson:"id" json:"id"`
	Name      string  `bson:"name" json:"name"`
	StageName *string `bson:"stageName" json:"stageName"`
}

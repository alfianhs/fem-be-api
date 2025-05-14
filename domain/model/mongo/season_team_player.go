package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SeasonTeamPlayer struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	SeasonTeam SeasonTeamFK       `bson:"seasonTeam" json:"seasonTeam"`
	Player     PlayerFK           `bson:"player" json:"player"`
	Position   string             `bson:"position" json:"position"`
	Image      MediaFK            `bson:"image" json:"image"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time         `bson:"deletedAt" json:"-"`
}

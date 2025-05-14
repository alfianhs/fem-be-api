package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SeasonTeam struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SeasonID  string             `bson:"seasonId" json:"seasonId"`
	Season    SeasonFK           `bson:"-" json:"season"`
	Team      TeamFK             `bson:"team" json:"team"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

type SeasonTeamFK struct {
	ID       string `bson:"id" json:"id"`
	SeasonID string `bson:"seasonId" json:"seasonId"`
	TeamID   string `bson:"teamId" json:"teamId"`
	Team     TeamFK `bson:"-" json:"team"`
}

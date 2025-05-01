package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Season struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Status       SeasonStatus       `bson:"status" json:"-"`
	StatusString string             `bson:"-" json:"status"`
	Logo         MediaFK            `bson:"logo" json:"logo"`
	Banner       MediaFK            `bson:"banner" json:"banner"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}

func (s *Season) Format() *Season {
	s.StatusString = SeasonStatusMap[s.Status].Name

	return s
}

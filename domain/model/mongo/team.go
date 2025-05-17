package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Logo       MediaFK            `bson:"logo" json:"logo"`
	IsSelected *bool              `bson:"-" json:"isSelected,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  *time.Time         `bson:"deletedAt" json:"-"`
}

type TeamFK struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
	Logo string `bson:"logo" json:"logo"`
}

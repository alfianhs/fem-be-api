package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Superadmin struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     *string            `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

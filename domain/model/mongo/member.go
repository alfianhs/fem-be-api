package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"-"`
	Phone         *string            `bson:"phone" json:"phone"`
	Age           *int               `bson:"age" json:"age"`
	Gender        *string            `bson:"gender" json:"gender"`
	EmailToken    string             `bson:"emailToken" json:"-"`
	PasswordToken string             `bson:"passwordToken" json:"-"`
	IsVerified    bool               `bson:"isVerified" json:"isVerified"`
	VerifiedAt    *time.Time         `bson:"verifiedAt" json:"-"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time         `bson:"deletedAt" json:"-"`
}

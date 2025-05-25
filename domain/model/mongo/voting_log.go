package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VotingLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Candidate CandidateFK        `bson:"candidate" json:"candidate"`
	MemberID  string             `bson:"memberId" json:"memberId"`
	Member    MemberFK           `bson:"-" json:"member"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

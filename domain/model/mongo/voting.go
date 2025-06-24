package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Voting struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	SeriesID         string             `bson:"seriesId" json:"seriesId"`
	Series           SeriesFK           `bson:"-" json:"series"`
	Title            string             `bson:"title" json:"title"`
	StartDate        time.Time          `bson:"startDate" json:"startDate"`
	EndDate          time.Time          `bson:"endDate" json:"endDate"`
	PerformancePoint PerformancePoint   `bson:"performancePoint" json:"performancePoint"`
	TotalVoter       int64              `bson:"totalVoter" json:"totalVoter"`
	Banner           MediaFK            `bson:"banner" json:"banner"`
	Status           VotingStatus       `bson:"status" json:"-"`
	StatusString     string             `bson:"-" json:"status"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time         `bson:"deletedAt" json:"-"`
}

type VotingFK struct {
	ID         string `bson:"id" json:"id"`
	Title      string `bson:"title" json:"title"`
	TotalVoter int64  `bson:"totalVoter" json:"totalVoter"`
}

type PerformancePoint struct {
	Goal   int64 `bson:"goal" json:"goal"`
	Assist int64 `bson:"assist" json:"assist"`
	Save   int64 `bson:"save" json:"save"`
}

func (v *Voting) Format() *Voting {
	v.StatusString = VotingStatusMap[v.Status].Name

	return v
}

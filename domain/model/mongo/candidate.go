package mongo_model

import (
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Candidate struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	VotingID         string             `bson:"votingId" json:"votingId"`
	Voting           VotingFK           `bson:"-" json:"voting"`
	SeasonTeam       SeasonTeamFK       `bson:"seasonTeam" json:"seasonTeam"`
	SeasonTeamPlayer SeasonTeamPlayerFK `bson:"seasonTeamPlayer" json:"seasonTeamPlayer"`
	Performance      Performance        `bson:"performance" json:"performance"`
	Voters           Voters             `bson:"voters" json:"voters"`
	IsChosen         *bool              `bson:"-" json:"isChosen,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt        *time.Time         `bson:"deletedAt" json:"-"`
}

type CandidateFK struct {
	ID       string `bson:"id" json:"id"`
	VotingID string `bson:"votingId" json:"votingId"`
}

type Voters struct {
	Count      int64   `bson:"count" json:"count"`
	Percentage float64 `bson:"-" json:"percentage"`
}

type Performance struct {
	TeamLeaderboard int64 `bson:"teamLeaderboard" json:"teamLeaderboard"`
	Goal            int64 `bson:"goal" json:"goal"`
	Assist          int64 `bson:"assist" json:"assist"`
	Save            int64 `bson:"save" json:"save"`
	Score           int64 `bson:"score" json:"score"`
}

func (c *Candidate) Format(v *VotingFK) *Candidate {
	if c.Voters.Count == 0 {
		c.Voters.Percentage = 0
		return c
	}
	percentage := (float64(c.Voters.Count) / float64(v.TotalVoter))
	c.Voters.Percentage = (math.Round(percentage*100*100) / 100)
	return c
}

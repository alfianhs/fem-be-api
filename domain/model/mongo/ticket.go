package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SeriesID  string             `bson:"seriesId" json:"seriesId"`
	Name      string             `bson:"name" json:"name"`
	Date      time.Time          `bson:"date" json:"date"`
	Price     float64            `bson:"price" json:"price"`
	Quota     TicketQuota        `bson:"quota" json:"quota"`
	Matchs    []TicketMatch      `bson:"matchs" json:"matchs"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

type TicketQuota struct {
	Stock     int64 `bson:"stock" json:"stock"`
	Used      int64 `bson:"used" json:"used"`
	Remaining int64 `bson:"-" json:"remaining"`
}

type TicketMatch struct {
	HomeSeasonTeamID string       `bson:"homeSeasonTeamId" json:"homeSeasonTeamId"`
	HomeSeasonTeam   SeasonTeamFK `bson:"-" json:"homeSeasonTeam"`
	AwaySeasonTeamID string       `bson:"awaySeasonTeamId" json:"awaySeasonTeamId"`
	AwaySeasonTeam   SeasonTeamFK `bson:"-" json:"awaySeasonTeam"`
	VenueID          string       `bson:"venueId" json:"venueId"`
	Venue            VenueFK      `bson:"-" json:"venue"`
	Time             string       `bson:"time" json:"time"`
}

func (t *Ticket) Format() *Ticket {
	t.Quota.Remaining = t.Quota.Stock - t.Quota.Used

	return t
}

package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Series struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	SeasonID     string             `bson:"seasonId" json:"seasonId"`
	Season       SeasonFK           `bson:"-" json:"season"`
	VenueID      string             `bson:"venueId" json:"venueId"`
	Venue        VenueFK            `bson:"-" json:"venue"`
	Name         string             `bson:"name" json:"name"`
	Price        float64            `bson:"price" json:"price"`
	MatchCount   int64              `bson:"matchCount" json:"matchCount"`
	StartDate    time.Time          `bson:"startDate" json:"startDate"`
	EndDate      time.Time          `bson:"endDate" json:"endDate"`
	Status       SeriesStatus       `bson:"status" json:"-"`
	StatusString string             `bson:"-" json:"status"`
	Tickets      []Ticket           `bson:"tickets" json:"tickets"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}

type SeriesFK struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func (s *Series) Format() *Series {
	s.StatusString = SeriesStatusMap[s.Status].Name

	return s
}

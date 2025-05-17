package request

import mongo_model "app/domain/model/mongo"

type SeriesCreateRequest struct {
	Name      string  `json:"name"`
	VenueID   string  `json:"venueId"`
	Price     float64 `json:"price"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
}

type SeriesUpdateRequest struct {
	Name      string                    `json:"name"`
	VenueID   string                    `json:"venueId"`
	Price     float64                   `json:"price"`
	StartDate string                    `json:"startDate"`
	EndDate   string                    `json:"endDate"`
	Status    *mongo_model.SeriesStatus `json:"status"`
}

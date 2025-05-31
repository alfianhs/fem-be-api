package request

import mongo_model "app/domain/model/mongo"

type VotingCreateRequest struct {
	SeriesID  string                   `form:"seriesId" validate:"required"`
	Title     string                   `form:"title" validate:"required"`
	StartDate string                   `form:"startDate" validate:"required"`
	EndDate   string                   `form:"endDate" validate:"required"`
	Status    mongo_model.VotingStatus `form:"status" validate:"required"`
}

type VotingUpdateRequest struct {
	SeriesID  string                   `form:"seriesId"`
	Title     string                   `form:"title"`
	StartDate string                   `form:"startDate"`
	EndDate   string                   `form:"endDate"`
	Status    mongo_model.VotingStatus `form:"status"`
}

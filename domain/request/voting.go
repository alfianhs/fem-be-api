package request

import mongo_model "app/domain/model/mongo"

type VotingCreateRequest struct {
	SeriesID  string                   `form:"seriesId" binding:"required"`
	Title     string                   `form:"title" binding:"required"`
	StartDate string                   `form:"startDate" binding:"required"`
	EndDate   string                   `form:"endDate" binding:"required"`
	Status    mongo_model.VotingStatus `form:"status" binding:"required"`
}

type VotingUpdateRequest struct {
	SeriesID  string                   `form:"seriesId"`
	Title     string                   `form:"title"`
	StartDate string                   `form:"startDate"`
	EndDate   string                   `form:"endDate"`
	Status    mongo_model.VotingStatus `form:"status"`
}

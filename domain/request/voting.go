package request

import mongo_model "app/domain/model/mongo"

type VotingCreateRequest struct {
	SeriesID    string                   `form:"seriesId" validate:"required"`
	Title       string                   `form:"title" validate:"required"`
	StartDate   string                   `form:"startDate" validate:"required"`
	EndDate     string                   `form:"endDate" validate:"required"`
	Status      mongo_model.VotingStatus `form:"status" validate:"required"`
	GoalPoint   int64                    `form:"goalPoint" validate:"required"`
	AssistPoint int64                    `form:"assistPoint" validate:"required"`
	SavePoint   int64                    `form:"savePoint" validate:"required"`
}

type VotingUpdateRequest struct {
	SeriesID    string                   `form:"seriesId"`
	Title       string                   `form:"title"`
	StartDate   string                   `form:"startDate"`
	EndDate     string                   `form:"endDate"`
	Status      mongo_model.VotingStatus `form:"status"`
	GoalPoint   int64                    `form:"goalPoint"`
	AssistPoint int64                    `form:"assistPoint"`
	SavePoint   int64                    `form:"savePoint"`
}

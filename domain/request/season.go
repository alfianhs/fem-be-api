package request

import mongo_model "app/domain/model/mongo"

type SeasonCreateRequest struct {
	Name string `form:"name" binding:"required"`
}

type SeasonUpdateRequest struct {
	Name string `form:"name" binding:"required"`
}

type SeasonStatusUpdateRequest struct {
	Status *mongo_model.SeasonStatus `json:"status"`
}

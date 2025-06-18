package member_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"
	"net/http"
)

func (u *memberAppUsecase) GetActiveSeasonDetail(ctx context.Context) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Active season not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, season.Format())
}

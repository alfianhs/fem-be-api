package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminAppUsecase) GetPlayerList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if queryParam.Get("search") != "" {
		fetchOptions["name"] = queryParam.Get("search")
	}

	// count total
	total := u.mongoDbRepo.CountPlayer(ctx, fetchOptions)
	if total == 0 {
		return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
			List:  []interface{}{},
			Limit: limit,
			Page:  page,
			Total: total,
		})
	}

	// sorting
	if queryParam.Get("sort") != "" {
		fetchOptions["sort"] = queryParam.Get("sort")
	}
	if queryParam.Get("dir") != "" {
		fetchOptions["dir"] = queryParam.Get("dir")
	}

	// fetch data
	cur, err := u.mongoDbRepo.FetchListPlayer(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var list []interface{}
	for cur.Next(ctx) {
		row := mongo_model.Player{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListPlayer Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		list = append(list, row)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		Limit: limit,
		Page:  page,
		Total: total,
		List:  list,
	})
}

func (u *superadminAppUsecase) GetPlayerDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	player, err := u.mongoDbRepo.FetchOnePlayer(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if player == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Player not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, player)
}

func (u *superadminAppUsecase) CreatePlayer(ctx context.Context, payload request.PlayerCreateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// create player
	now := time.Now()
	player := mongo_model.Player{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save
	err := u.mongoDbRepo.CreateOnePlayer(ctx, &player)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Create player success", nil, player)
}

func (u *superadminAppUsecase) UpdatePlayer(ctx context.Context, id string, payload request.PlayerUpdateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// get player
	player, err := u.mongoDbRepo.FetchOnePlayer(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if player == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Player not found", nil, nil)
	}

	// update player
	player.Name = payload.Name
	player.UpdatedAt = time.Now()

	// save
	err = u.mongoDbRepo.UpdatePartialPlayer(ctx, map[string]interface{}{
		"id": player.ID,
	}, map[string]interface{}{
		"name":      player.Name,
		"updatedAt": player.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// update season team player in bg
	go u.updateActiveSeasonTeamPlayerBackground(context.Background(), player.ID.Hex(), player)

	return helpers.NewResponse(http.StatusOK, "Update player success", nil, player)
}

func (u *superadminAppUsecase) DeletePlayer(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get player
	player, err := u.mongoDbRepo.FetchOnePlayer(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if player == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Player not found", nil, nil)
	}

	// check if player is in active season team player
	activeSeason, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if activeSeason != nil {
		seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{
			"seasonId": activeSeason.ID.Hex(),
			"playerId": player.ID.Hex(),
		})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		if seasonTeamPlayer != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Player is in active season team player", nil, nil)
		}
	}

	// delete player
	now := time.Now()
	player.DeletedAt = &now

	// save
	err = u.mongoDbRepo.UpdatePartialPlayer(ctx, map[string]interface{}{
		"id": player.ID,
	}, map[string]interface{}{
		"deletedAt": player.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Delete player success", nil, nil)
}

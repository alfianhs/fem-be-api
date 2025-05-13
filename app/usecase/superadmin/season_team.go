package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminAppUsecase) GetSeasonTeamsList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if queryParam.Get("seasonId") != "" {
		fetchOptions["seasonId"] = queryParam.Get("seasonId")
	}

	// count total
	total := u.mongoDbRepo.CountSeasonTeam(ctx, fetchOptions)
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

	// fetch list
	cur, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var seasonTeams []mongo_model.SeasonTeam
	for cur.Next(ctx) {
		row := mongo_model.SeasonTeam{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListSeasonTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		seasonTeams = append(seasonTeams, row)
	}

	// get seasons ids
	seasonIds := helpers.ExtractIds(seasonTeams, func(m mongo_model.SeasonTeam) string {
		return m.SeasonID
	})
	seasonFetchOptions := map[string]interface{}{
		"ids": seasonIds,
	}

	// fetch seasons
	seasons, err := u.mongoDbRepo.FetchListSeason(ctx, seasonFetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer seasons.Close(ctx)

	seasonMap := make(map[string]mongo_model.Season)
	for seasons.Next(ctx) {
		var season mongo_model.Season
		if err := seasons.Decode(&season); err != nil {
			logrus.Error("Season Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		seasonMap[season.ID.Hex()] = season
	}

	var list []interface{}
	for _, row := range seasonTeams {
		season, ok := seasonMap[row.SeasonID]
		if ok {
			row.Season = mongo_model.SeasonFK{
				ID:   row.SeasonID,
				Name: season.Name,
			}
		} else {
			row.Season = mongo_model.SeasonFK{
				ID:   row.SeasonID,
				Name: "",
			}
			logrus.Error("Season with id ", row.SeasonID, " not found")
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

func (u *superadminAppUsecase) GetSeasonTeamDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeam == nil {
		return helpers.NewResponse(http.StatusBadRequest, "SeasonTeam not found", nil, nil)
	}

	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": seasonTeam.SeasonID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		season = &mongo_model.Season{}
		logrus.Error("Season with id ", seasonTeam.SeasonID, " not found")
	}

	seasonTeam.Season = mongo_model.SeasonFK{
		ID:   season.ID.Hex(),
		Name: season.Name,
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, seasonTeam)
}

func (u *superadminAppUsecase) CreateSeasonTeam(ctx context.Context, payload request.SeasonTeamCreateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if len(payload.TeamIds) == 0 {
		errValidation["teamIds"] = "TeamIds field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// get active season
	activeSeason, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if activeSeason == nil {
		return helpers.NewResponse(http.StatusBadRequest, "There is no active season", nil, nil)
	}

	// check if existing
	existing, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, map[string]interface{}{
		"seasonId": activeSeason.ID.Hex(),
		"team.ids": payload.TeamIds,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer existing.Close(ctx)
	validationMap := make(map[string]string)
	for existing.Next(ctx) {
		var seasonTeam mongo_model.SeasonTeam
		if err := existing.Decode(&seasonTeam); err != nil {
			logrus.Error("SeasonTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		validationMap["teamIds"] = fmt.Sprintf("Team with id %s already exists", seasonTeam.Team.ID)
	}
	if len(validationMap) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", validationMap, nil)
	}

	// get teams
	teamFetchOptions := map[string]interface{}{
		"ids": payload.TeamIds,
	}
	teams, err := u.mongoDbRepo.FetchListTeam(ctx, teamFetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer teams.Close(ctx)

	teamMap := make(map[string]mongo_model.Team)
	for teams.Next(ctx) {
		var team mongo_model.Team
		if err := teams.Decode(&team); err != nil {
			logrus.Error("Team Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		teamMap[team.ID.Hex()] = team
	}

	now := time.Now()

	var seasonTeams []*mongo_model.SeasonTeam
	for _, teamId := range payload.TeamIds {
		team, ok := teamMap[teamId]
		if ok {
			seasonTeams = append(seasonTeams, &mongo_model.SeasonTeam{
				ID:       primitive.NewObjectID(),
				SeasonID: activeSeason.ID.Hex(),
				Season: mongo_model.SeasonFK{
					ID:   activeSeason.ID.Hex(),
					Name: activeSeason.Name,
				},
				Team: mongo_model.TeamFK{
					ID:   team.ID.Hex(),
					Name: team.Name,
					Logo: team.Logo.URL,
				},
				CreatedAt: now,
				UpdatedAt: now,
			})
		} else {
			return helpers.NewResponse(http.StatusBadRequest, "Team with id "+teamId+" not found", nil, nil)
		}
	}

	// save
	err = u.mongoDbRepo.CreateManySeasonTeam(ctx, seasonTeams)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Success", nil, seasonTeams)
}

func (u *superadminAppUsecase) DeleteSeasonTeam(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeam == nil {
		return helpers.NewResponse(http.StatusBadRequest, "SeasonTeam not found", nil, nil)
	}

	now := time.Now()
	seasonTeam.DeletedAt = &now

	err = u.mongoDbRepo.UpdatePartialSeasonTeam(ctx, map[string]interface{}{
		"id": id,
	}, map[string]interface{}{
		"deletedAt": seasonTeam.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, nil)
}

func (u *superadminAppUsecase) updateActiveSeasonTeamBackground(ctx context.Context, teamId string, team *mongo_model.Team) {
	// get active season
	activeSeason, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		logrus.Error("get active season background error :", err)
	}
	if activeSeason == nil {
		return
	}

	// update
	err = u.mongoDbRepo.UpdatePartialSeasonTeam(ctx, map[string]interface{}{
		"seasonId": activeSeason.ID.Hex(),
		"team.id":  teamId,
	}, map[string]interface{}{
		"team.name": team.Name,
		"team.logo": team.Logo.URL,
	})
	if err != nil {
		logrus.Error("update season team background error :", err)
	}
}

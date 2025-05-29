package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminAppUsecase) GetTeamsList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// handle selected in season team
	seasonIdParam := queryParam.Get("seasonId")
	var selectedTeamMap map[string]struct{}
	if seasonIdParam != "" {
		var seasonId string
		if seasonIdParam == "active" {
			season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
				"status": mongo_model.SeasonStatusActive,
			})
			if err != nil {
				return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
			}
			if season == nil {
				return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
					List:  []interface{}{},
					Limit: limit,
					Page:  page,
					Total: 0,
				})
			}
			seasonId = season.ID.Hex()
		} else {
			seasonId = seasonIdParam
		}

		// fetch season teams
		seasonTeamsCursor, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, map[string]interface{}{
			"seasonId": seasonId,
		})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		defer seasonTeamsCursor.Close(ctx)

		selectedTeamMap = make(map[string]struct{})
		for seasonTeamsCursor.Next(ctx) {
			var st mongo_model.SeasonTeam
			err = seasonTeamsCursor.Decode(&st)
			if err != nil {
				logrus.Error("SeasonTeam Decode:", err)
				return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
			}
			selectedTeamMap[st.Team.ID] = struct{}{}
		}
	}

	// filtering
	if queryParam.Get("search") != "" {
		fetchOptions["name"] = queryParam.Get("search")
	}

	// count total
	total := u.mongoDbRepo.CountTeam(ctx, fetchOptions)
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
	cur, err := u.mongoDbRepo.FetchListTeam(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var list []interface{}
	for cur.Next(ctx) {
		row := mongo_model.Team{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		// handle selected in season team
		var isSelected *bool
		if selectedTeamMap != nil {
			selected := false
			if _, found := selectedTeamMap[row.ID.Hex()]; found {
				selected = true
			}
			isSelected = &selected
		}
		row.IsSelected = isSelected

		list = append(list, row)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		Limit: limit,
		Page:  page,
		Total: total,
		List:  list,
	})
}

func (u *superadminAppUsecase) GetTeamDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	team, err := u.mongoDbRepo.FetchOneTeam(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if team == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Team not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, team)
}

func (u *superadminAppUsecase) CreateTeam(ctx context.Context, payload request.TeamCreateRequest, request *http.Request) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	logoFile, logoFileHeader, err := request.FormFile("logo")
	if err != nil {
		errValidation["logo"] = "Logo field is required"
	} else {
		err := helpers.ImageUploadValidation(logoFile, logoFileHeader)
		if err != nil {
			errValidation["logo"] = err.Error()
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// upload images to s3
	logoObjName := "team/logo/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(logoFileHeader.Filename)
	logoImgType := logoFileHeader.Header.Get("Content-Type")
	uploadLogo, err := u.s3Repo.UploadFilePublic(ctx, logoObjName, logoFile, logoImgType)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
	}

	// create media
	media := &mongo_model.Media{
		ID:          primitive.NewObjectID(),
		Name:        logoFileHeader.Filename,
		Provider:    "s3",
		ProviderKey: uploadLogo.Key,
		Type:        logoImgType,
		Size:        logoFileHeader.Size,
		URL:         uploadLogo.URL,
		IsUsed:      true,
		IsPrivate:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// save media
	err = u.mongoDbRepo.CreateOneMedia(ctx, media)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// create team
	team := mongo_model.Team{
		ID:   primitive.NewObjectID(),
		Name: payload.Name,
		Logo: mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save team
	err = u.mongoDbRepo.CreateOneTeam(ctx, &team)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Create team success", nil, team)
}

func (u *superadminAppUsecase) UpdateTeam(ctx context.Context, id string, payload request.TeamUpdateRequest, request *http.Request) helpers.Response {
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

	// check team
	team, err := u.mongoDbRepo.FetchOneTeam(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if team == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Team not found", nil, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// update file if upload exist
	logoFile, logoFileHeader, err := request.FormFile("logo")
	if err == nil {
		// upload new media
		logoObjName := "team/logo/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(logoFileHeader.Filename)
		logoImgType := logoFileHeader.Header.Get("Content-Type")
		uploadLogo, err := u.s3Repo.UploadFilePublic(ctx, logoObjName, logoFile, logoImgType)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
		}

		// create new media
		media := mongo_model.Media{
			ID:          primitive.NewObjectID(),
			Name:        logoFileHeader.Filename,
			Provider:    "s3",
			ProviderKey: uploadLogo.Key,
			Type:        logoImgType,
			Size:        logoFileHeader.Size,
			URL:         uploadLogo.URL,
			IsUsed:      true,
			IsPrivate:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// save media
		err = u.mongoDbRepo.CreateOneMedia(ctx, &media)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		// update old media in bg
		go u.markMediaAsUnusedByIds(context.Background(), []string{team.Logo.ID})

		// update team
		team.Logo = mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		}
	}

	// update team
	team.Name = payload.Name
	team.UpdatedAt = now

	// save team
	err = u.mongoDbRepo.UpdatePartialTeam(ctx, map[string]interface{}{
		"id": team.ID,
	}, map[string]interface{}{
		"name":      team.Name,
		"logo":      team.Logo,
		"updatedAt": team.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// bg update active season team
	go u.updateActiveSeasonTeamBackground(context.Background(), team.ID.Hex(), team)

	return helpers.NewResponse(http.StatusOK, "Update team success", nil, team)
}

func (u *superadminAppUsecase) DeleteTeam(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get team
	team, err := u.mongoDbRepo.FetchOneTeam(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if team == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Team not found", nil, nil)
	}

	// check if team is in active season team
	activeSeason, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if activeSeason != nil {
		seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{
			"season.id": activeSeason.ID.Hex(),
			"team.id":   team.ID.Hex(),
		})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		if seasonTeam != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Team is in active season", nil, nil)
		}
	}

	// delete team
	now := time.Now()
	team.DeletedAt = &now

	// update medias
	go u.markMediaAsUnusedByIds(context.Background(), []string{team.Logo.ID})

	// save
	err = u.mongoDbRepo.UpdatePartialTeam(ctx, map[string]interface{}{
		"id": team.ID,
	}, map[string]interface{}{
		"deletedAt": team.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Delete team success", nil, nil)
}

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

func (u *superadminAppUsecase) GetSeasonTeamPlayersList(ctx context.Context, queryParam url.Values) helpers.Response {
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
	if queryParam.Get("seasonTeamId") != "" {
		fetchOptions["seasonTeamId"] = queryParam.Get("seasonTeamId")
	}

	// count total
	total := u.mongoDbRepo.CountSeasonTeamPlayer(ctx, fetchOptions)
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
	cur, err := u.mongoDbRepo.FetchListSeasonTeamPlayer(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var seasonTeamPlayers []mongo_model.SeasonTeamPlayer
	for cur.Next(ctx) {
		row := mongo_model.SeasonTeamPlayer{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("SeasonTeamPlayer Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		seasonTeamPlayers = append(seasonTeamPlayers, row)
	}

	teamIds := helpers.ExtractIds(seasonTeamPlayers, func(m mongo_model.SeasonTeamPlayer) string {
		return m.SeasonTeam.TeamID
	})
	teamFetchOptions := map[string]interface{}{
		"ids": teamIds,
	}

	// fetch teams
	teamCur, err := u.mongoDbRepo.FetchListTeam(ctx, teamFetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer teamCur.Close(ctx)

	teamMap := make(map[string]mongo_model.Team)
	for teamCur.Next(ctx) {
		row := mongo_model.Team{}
		err := teamCur.Decode(&row)
		if err != nil {
			logrus.Error("Team Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		teamMap[row.ID.Hex()] = row
	}

	var list []interface{}
	for _, row := range seasonTeamPlayers {
		team, ok := teamMap[row.SeasonTeam.TeamID]
		if ok {
			row.SeasonTeam.Team = mongo_model.TeamFK{
				ID:   team.ID.Hex(),
				Name: team.Name,
				Logo: team.Logo.URL,
			}
		} else {
			logrus.Error("team with id ", row.SeasonTeam.TeamID, " not found")
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

func (u *superadminAppUsecase) GetSeasonTeamPlayerDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeamPlayer == nil {
		return helpers.NewResponse(http.StatusBadRequest, "SeasonTeamPlayer not found", nil, nil)
	}

	team, err := u.mongoDbRepo.FetchOneTeam(ctx, map[string]interface{}{
		"id": seasonTeamPlayer.SeasonTeam.TeamID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if team == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Team not found", nil, nil)
	}
	seasonTeamPlayer.SeasonTeam.Team = mongo_model.TeamFK{
		ID:   team.ID.Hex(),
		Name: team.Name,
		Logo: team.Logo.URL,
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, seasonTeamPlayer)
}

func (u *superadminAppUsecase) GetPlayerPositionsList(ctx context.Context) helpers.Response {
	return helpers.NewResponse(http.StatusOK, "Success", nil, mongo_model.PlayerPositionList)
}

func (u *superadminAppUsecase) CreateSeasonTeamPlayer(ctx context.Context, payload request.SeasonTeamPlayerCreateRequest, request *http.Request) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.SeasonTeamID == "" {
		errValidation["seasonTeamId"] = "SeasonTeamId field is required"
	}
	if payload.PlayerID == "" {
		errValidation["playerId"] = "PlayerId field is required"
	}
	if payload.Position == "" {
		errValidation["position"] = "Position field is required"
	} else {
		if !helpers.InArrayString(mongo_model.PlayerPositionList, payload.Position) {
			errValidation["position"] = "Invalid position, required position: " + helpers.ArrayStringtoString(mongo_model.PlayerPositionList)
		}
	}
	imageFile, imageFileHeader, err := request.FormFile("image")
	if err != nil {
		errValidation["image"] = "Image field is required"
	} else {
		err := helpers.ImageUploadValidation(imageFile, imageFileHeader)
		if err != nil {
			errValidation["image"] = err.Error()
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check season team
	seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{
		"id": payload.SeasonTeamID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeam == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season team not found", nil, nil)
	}

	// check player
	player, err := u.mongoDbRepo.FetchOnePlayer(ctx, map[string]interface{}{
		"id": payload.PlayerID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if player == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Player not found", nil, nil)
	}

	// check existing
	existing, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{
		"seasonId": seasonTeam.SeasonID,
		"playerId": player.ID.Hex(),
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if existing != nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season team player already exists in this season", nil, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// upload images to s3
	imageObjName := "player/images/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(imageFileHeader.Filename)
	imageType := imageFileHeader.Header.Get("Content-Type")
	uploadImage, err := u.s3Repo.UploadFilePublic(ctx, imageObjName, imageFile, imageType)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
	}

	// create media
	media := &mongo_model.Media{
		ID:          primitive.NewObjectID(),
		Name:        imageFileHeader.Filename,
		Provider:    "s3",
		ProviderKey: uploadImage.Key,
		Type:        imageType,
		Size:        imageFileHeader.Size,
		URL:         uploadImage.URL,
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

	// create season team player
	seasonTeamPlayer := &mongo_model.SeasonTeamPlayer{
		ID: primitive.NewObjectID(),
		SeasonTeam: mongo_model.SeasonTeamFK{
			ID:       seasonTeam.ID.Hex(),
			SeasonID: seasonTeam.SeasonID,
			TeamID:   seasonTeam.Team.ID,
			Team: mongo_model.TeamFK{
				ID:   seasonTeam.Team.ID,
				Name: seasonTeam.Team.Name,
				Logo: seasonTeam.Team.Logo,
			},
		},
		Player: mongo_model.PlayerFK{
			ID:   player.ID.Hex(),
			Name: player.Name,
		},
		Position: payload.Position,
		Image: mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			ProviderKey: media.ProviderKey,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save
	err = u.mongoDbRepo.CreateOneSeasonTeamPlayer(ctx, seasonTeamPlayer)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Success", nil, seasonTeamPlayer)
}

func (u *superadminAppUsecase) UpdateSeasonTeamPlayer(ctx context.Context, id string, payload request.SeasonTeamPlayerUpdateRequest, request *http.Request) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Position == "" {
		errValidation["position"] = "Position field is required"
	} else {
		if !helpers.InArrayString(mongo_model.PlayerPositionList, payload.Position) {
			errValidation["position"] = "Invalid position, required position: " + helpers.ArrayStringtoString(mongo_model.PlayerPositionList)
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check season team player
	seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeamPlayer == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season team player not found", nil, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// update file if upload exist
	imageFile, imageFileHeader, err := request.FormFile("image")
	if err == nil {
		// upload new media
		imageObjName := "player/images/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(imageFileHeader.Filename)
		imageImgType := imageFileHeader.Header.Get("Content-Type")
		uploadImage, err := u.s3Repo.UploadFilePublic(ctx, imageObjName, imageFile, imageImgType)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
		}

		// create new media
		media := mongo_model.Media{
			ID:          primitive.NewObjectID(),
			Name:        imageFileHeader.Filename,
			Provider:    "s3",
			ProviderKey: uploadImage.Key,
			Type:        imageImgType,
			Size:        imageFileHeader.Size,
			URL:         uploadImage.URL,
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
		go u.markMediaAsUnusedByIds(context.Background(), []string{seasonTeamPlayer.Image.ID})

		// update seasonTeamPlayer
		seasonTeamPlayer.Image = mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		}
	}

	// update seasonTeamPlayer
	seasonTeamPlayer.Position = payload.Position
	seasonTeamPlayer.UpdatedAt = now

	// save
	err = u.mongoDbRepo.UpdatePartialSeasonTeamPlayer(ctx, map[string]interface{}{
		"id": seasonTeamPlayer.ID,
	}, map[string]interface{}{
		"position":  seasonTeamPlayer.Position,
		"image":     seasonTeamPlayer.Image,
		"updatedAt": seasonTeamPlayer.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, seasonTeamPlayer)
}

func (u *superadminAppUsecase) DeleteSeasonTeamPlayer(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check season team player
	seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeamPlayer == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season team player not found", nil, nil)
	}

	// delete season team player
	now := time.Now()
	seasonTeamPlayer.DeletedAt = &now
	err = u.mongoDbRepo.UpdatePartialSeasonTeamPlayer(ctx, map[string]interface{}{
		"id": seasonTeamPlayer.ID,
	}, map[string]interface{}{
		"deletedAt": seasonTeamPlayer.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, nil)
}

func (u *superadminAppUsecase) updateActiveSeasonTeamPlayerBackground(ctx context.Context, player *mongo_model.Player) {
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
	err = u.mongoDbRepo.UpdatePartialSeasonTeamPlayer(ctx, map[string]interface{}{
		"seasonId": activeSeason.ID.Hex(),
		"playerId": player.ID.Hex(),
	}, map[string]interface{}{
		"player.name":      player.Name,
		"player.stageName": player.StageName,
	})
	if err != nil {
		logrus.Error("update season team player background error :", err)
	}
}

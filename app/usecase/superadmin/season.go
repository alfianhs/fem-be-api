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

func (u *superadminAppUsecase) GetSeasonsList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// count total
	total := u.mongoDbRepo.CountSeason(ctx, fetchOptions)
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
	cur, err := u.mongoDbRepo.FetchListSeason(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var list []interface{}
	for cur.Next(ctx) {
		row := mongo_model.Season{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListSeason Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		list = append(list, row.Format())
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		Limit: limit,
		Page:  page,
		Total: total,
		List:  list,
	})
}

func (u *superadminAppUsecase) GetSeasonDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, season.Format())
}

func (u *superadminAppUsecase) CreateSeason(ctx context.Context, payload request.SeasonCreateRequest, request *http.Request) helpers.Response {
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
	bannerFile, bannerFileHeader, err := request.FormFile("banner")
	if err != nil {
		errValidation["banner"] = "Logo field is required"
	} else {
		err := helpers.ImageUploadValidation(bannerFile, bannerFileHeader)
		if err != nil {
			errValidation["banner"] = err.Error()
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// upload images to s3
	logoObjName := "season/logo/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(logoFileHeader.Filename)
	logoImgType := logoFileHeader.Header.Get("Content-Type")
	uploadLogo, err := u.s3Repo.UploadFilePublic(ctx, logoObjName, logoFile, logoImgType)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
	}
	bannerObjName := "season/banner/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(bannerFileHeader.Filename)
	bannerImgType := bannerFileHeader.Header.Get("Content-Type")
	uploadBanner, err := u.s3Repo.UploadFilePublic(ctx, bannerObjName, bannerFile, bannerImgType)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
	}

	// create media
	medias := []*mongo_model.Media{
		{
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
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        bannerFileHeader.Filename,
			Provider:    "s3",
			ProviderKey: uploadBanner.Key,
			Type:        bannerImgType,
			Size:        bannerFileHeader.Size,
			URL:         uploadBanner.URL,
			IsUsed:      true,
			IsPrivate:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// save many medias
	err = u.mongoDbRepo.CreateManyMedia(ctx, medias)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// create season
	season := mongo_model.Season{
		ID:     primitive.NewObjectID(),
		Name:   payload.Name,
		Status: mongo_model.SeasonStatusInactive,
		Logo: mongo_model.MediaFK{
			ID:          medias[0].ID.Hex(),
			Name:        medias[0].Name,
			Size:        medias[0].Size,
			URL:         medias[0].URL,
			Type:        medias[0].Type,
			IsPrivate:   medias[0].IsPrivate,
			ProviderKey: medias[0].ProviderKey,
		},
		Banner: mongo_model.MediaFK{
			ID:          medias[1].ID.Hex(),
			Name:        medias[1].Name,
			Size:        medias[1].Size,
			URL:         medias[1].URL,
			Type:        medias[1].Type,
			IsPrivate:   medias[1].IsPrivate,
			ProviderKey: medias[1].ProviderKey,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save season
	err = u.mongoDbRepo.CreateOneSeason(ctx, &season)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Create season success", nil, season.Format())
}

func (u *superadminAppUsecase) UpdateSeason(ctx context.Context, id string, payload request.SeasonUpdateRequest, requestHttp *http.Request) helpers.Response {
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

	// check season
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season not found", nil, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// update file if upload exist
	logoFile, logoFileHeader, err := requestHttp.FormFile("logo")
	if err == nil {
		// upload new media
		logoObjName := "season/logo/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(logoFileHeader.Filename)
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
		go u.markMediaAsUnusedByIds(context.Background(), []string{season.Logo.ID})

		// update season
		season.Logo = mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		}
	}
	bannerFile, bannerFileHeader, err := requestHttp.FormFile("banner")
	if err == nil {
		// upload new media
		bannerObjName := "season/banner/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(bannerFileHeader.Filename)
		bannerImgType := bannerFileHeader.Header.Get("Content-Type")
		uploadBanner, err := u.s3Repo.UploadFilePublic(ctx, bannerObjName, bannerFile, bannerImgType)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
		}

		// create new media
		media := mongo_model.Media{
			ID:          primitive.NewObjectID(),
			Name:        bannerFileHeader.Filename,
			Provider:    "s3",
			ProviderKey: uploadBanner.Key,
			Type:        bannerImgType,
			Size:        bannerFileHeader.Size,
			URL:         uploadBanner.URL,
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
		go u.markMediaAsUnusedByIds(context.Background(), []string{season.Banner.ID})

		// update season
		season.Banner = mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		}
	}

	// update season
	season.Name = payload.Name
	season.UpdatedAt = now

	// save season
	err = u.mongoDbRepo.UpdatePartialSeason(ctx, map[string]interface{}{
		"id": season.ID,
	}, map[string]interface{}{
		"name":      season.Name,
		"logo":      season.Logo,
		"banner":    season.Banner,
		"updatedAt": season.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Update season success", nil, season.Format())
}

func (u *superadminAppUsecase) DeleteSeason(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get season
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season not found", nil, nil)
	}

	// delete season
	now := time.Now()
	season.DeletedAt = &now

	// update medias
	mediaIds := []string{
		season.Logo.ID,
		season.Banner.ID,
	}
	go u.markMediaAsUnusedByIds(context.Background(), mediaIds)

	// save
	err = u.mongoDbRepo.UpdatePartialSeason(ctx, map[string]interface{}{
		"id": season.ID,
	}, map[string]interface{}{
		"deletedAt": season.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Delete season success", nil, nil)
}

func (u *superadminAppUsecase) UpdateSeasonStatus(ctx context.Context, id string, payload request.SeasonStatusUpdateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Status == nil {
		errValidation["status"] = "Status field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusBadRequest, "Invalid payload", errValidation, nil)
	}

	// get season
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season not found", nil, nil)
	}

	// get active season
	activeSeason, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if activeSeason != nil && activeSeason.ID != season.ID && *payload.Status == mongo_model.SeasonStatusActive {
		return helpers.NewResponse(http.StatusBadRequest, "Cannot activate: another season is already active", nil, nil)
	}

	// update season
	now := time.Now()
	season.Status = *payload.Status
	season.UpdatedAt = now

	// save season
	err = u.mongoDbRepo.UpdatePartialSeason(ctx, map[string]interface{}{
		"id": season.ID,
	}, map[string]interface{}{
		"status":    season.Status,
		"updatedAt": season.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Update season status success", nil, season.Format())
}

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

func (u *superadminAppUsecase) GetVotingList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, offset, limit := helpers.GetOffsetLimit(queryParam)
	opts := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	// filtering
	if s := queryParam.Get("status"); s != "" {
		statusInt, err := strconv.Atoi(s)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Invalid status", nil, nil)
		}
		opts["status"] = mongo_model.VotingStatus(statusInt)
	}
	if queryParam.Get("search") != "" {
		opts["title"] = queryParam.Get("search")
	}

	total := u.mongoDbRepo.CountVoting(ctx, opts)
	if total == 0 {
		return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
			List:  []interface{}{},
			Limit: limit,
			Page:  page,
			Total: total,
		})
	}

	// sorting
	if sort := queryParam.Get("sort"); sort != "" {
		opts["sort"] = sort
	}
	if dir := queryParam.Get("dir"); dir != "" {
		opts["dir"] = dir
	}

	// fetch data
	cur, err := u.mongoDbRepo.FetchListVoting(ctx, opts)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var votings []mongo_model.Voting
	for cur.Next(ctx) {
		var v mongo_model.Voting
		if err := cur.Decode(&v); err != nil {
			logrus.Error("Voting Decode Error:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		votings = append(votings, v)
	}

	// extract series ids
	seriesIds := helpers.ExtractIds(votings, func(m mongo_model.Voting) string {
		return m.SeriesID
	})

	// fetch series data
	series, err := u.mongoDbRepo.FetchListSeries(ctx, map[string]interface{}{
		"ids": seriesIds,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer series.Close(ctx)

	seriesMap := make(map[string]mongo_model.Series)
	for series.Next(ctx) {
		var s mongo_model.Series
		if err := series.Decode(&s); err != nil {
			logrus.Error("Series Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		seriesMap[s.ID.Hex()] = s
	}

	// format voting data
	var list []interface{}
	for _, row := range votings {
		s, ok := seriesMap[row.SeriesID]
		if ok {
			row.Series = mongo_model.SeriesFK{
				ID:   s.ID.Hex(),
				Name: s.Name,
			}
		} else {
			row.Series = mongo_model.SeriesFK{
				ID:   row.SeriesID,
				Name: "",
			}

			logrus.Error("Series with id ", row.SeriesID, " not found")
		}

		list = append(list, row.Format())
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		List:  list,
		Limit: limit,
		Page:  page,
		Total: total,
	})
}

func (u *superadminAppUsecase) GetVotingDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting == nil {
		return helpers.NewResponse(http.StatusNotFound, "Voting not found", nil, nil)
	}

	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": voting.SeriesID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		series = &mongo_model.Series{}
		logrus.Error("Series with id ", voting.SeriesID, " not found")
	}

	voting.Series = mongo_model.SeriesFK{
		ID:   series.ID.Hex(),
		Name: series.Name,
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, voting.Format())
}

func (u *superadminAppUsecase) CreateVoting(ctx context.Context, payload request.VotingCreateRequest, request *http.Request) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1) Validate payload
	errs := map[string]string{}
	if payload.Title == "" {
		errs["title"] = "Title is required"
	}
	if payload.SeriesID == "" {
		errs["seriesId"] = "Series ID is required"
	}
	if payload.StartDate == "" {
		errs["startDate"] = "Start date is required"
	}
	if payload.EndDate == "" {
		errs["endDate"] = "End date is required"
	}
	if payload.Status == 0 {
		errs["status"] = "Status is required"
	}
	if payload.GoalPoint == 0 {
		errs["goalPoint"] = "Goal point is required"
	}
	if payload.AssistPoint == 0 {
		errs["assistPoint"] = "Assist point is required"
	}
	if payload.SavePoint == 0 {
		errs["savePoint"] = "Save point is required"
	}
	bannerFile, bannerFileHeader, err := request.FormFile("banner")
	if err != nil {
		errs["banner"] = "Banner field is required"
	} else {
		err := helpers.ImageUploadValidation(bannerFile, bannerFileHeader)
		if err != nil {
			errs["banner"] = err.Error()
		}
	}
	if payload.Status == mongo_model.VotingStatusActive {
		errs["status"] = "Status cannot be active in create data"
	}
	if len(errs) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation error", errs, nil)
	}

	now := time.Now()
	year, month, _ := now.Date()

	// upload banner to s3
	bannerObjName := "voting/banner/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(bannerFileHeader.Filename)
	bannerImgType := bannerFileHeader.Header.Get("Content-Type")
	uploadBanner, err := u.s3Repo.UploadFilePublic(ctx, bannerObjName, bannerFile, bannerImgType)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Something went wrong", nil, nil)
	}

	// create media
	media := &mongo_model.Media{
		ID:          primitive.NewObjectID(),
		Name:        bannerFileHeader.Filename,
		Type:        bannerImgType,
		Size:        bannerFileHeader.Size,
		Provider:    "s3",
		ProviderKey: bannerObjName,
		URL:         uploadBanner.URL,
		IsUsed:      true,
		IsPrivate:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// insert media
	err = u.mongoDbRepo.CreateOneMedia(ctx, media)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// 2) Parse dates
	start, err := time.Parse(time.RFC3339, payload.StartDate)
	if err != nil {
		errs["startDate"] = "Start date format is invalid"
	}
	end, err2 := time.Parse(time.RFC3339, payload.EndDate)
	if err2 != nil {
		errs["endDate"] = "End date format is invalid"
	}

	// validate range date
	if start.After(end) || start.Format("2006-01-02") == end.Format("2006-01-02") {
		return helpers.NewResponse(http.StatusBadRequest, "Start date must be before end date", nil, nil)
	}
	if len(errs) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation error", errs, nil)
	}

	// 3) Ensure series exists
	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{"id": payload.SeriesID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
	}

	// 4) Build model & insert

	voting := &mongo_model.Voting{
		ID:         primitive.NewObjectID(),
		SeriesID:   series.ID.Hex(),
		Title:      payload.Title,
		StartDate:  helpers.SetToStartOfDayWIB(start),
		EndDate:    helpers.SetToEndOfDayWIB(end),
		TotalVoter: 0,
		Status:     mongo_model.VotingStatus(payload.Status),
		PerformancePoint: mongo_model.PerformancePoint{
			Goal:   payload.GoalPoint,
			Assist: payload.AssistPoint,
			Save:   payload.SavePoint,
		},
		CreatedAt: now,
		UpdatedAt: now,
		Banner: mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   false,
			ProviderKey: media.ProviderKey,
		},
	}

	if err := u.mongoDbRepo.CreateOneVoting(ctx, voting); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Voting created", nil, voting.Format())
}

func (u *superadminAppUsecase) UpdateVoting(ctx context.Context, id string, payload request.VotingUpdateRequest, request *http.Request) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check voting
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Voting not found", nil, nil)
	}

	// update if not empty
	if payload.SeriesID != "" {
		series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
			"id": payload.SeriesID,
		})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		if series == nil {
			return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
		}
		voting.SeriesID = series.ID.Hex()
	}
	if payload.Title != "" {
		voting.Title = payload.Title
	}
	if payload.StartDate != "" {
		start, err := time.Parse(time.RFC3339, payload.StartDate)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Start date format is invalid", nil, nil)
		}
		voting.StartDate = helpers.SetToStartOfDayWIB(start)
	}
	if payload.EndDate != "" {
		end, err := time.Parse(time.RFC3339, payload.EndDate)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "End date format is invalid", nil, nil)
		}
		voting.EndDate = helpers.SetToEndOfDayWIB(end)
	}
	if payload.Status != 0 {
		voting.Status = mongo_model.VotingStatus(payload.Status)
	}
	if payload.GoalPoint != 0 {
		voting.PerformancePoint.Goal = payload.GoalPoint
	}
	if payload.AssistPoint != 0 {
		voting.PerformancePoint.Assist = payload.AssistPoint
	}
	if payload.SavePoint != 0 {
		voting.PerformancePoint.Save = payload.SavePoint
	}

	// set timestamp
	now := time.Now()
	year, month, _ := now.Date()

	// update file if upload exist
	bannerFile, bannerFileHeader, err := request.FormFile("banner")
	if err == nil {
		// upload new media
		bannerObjName := "voting/banner/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(bannerFileHeader.Filename)
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
		go u.markMediaAsUnusedByIds(context.Background(), []string{voting.Banner.ID})

		// update voting
		voting.Banner = mongo_model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		}
	}

	// update voting
	voting.UpdatedAt = now

	err = u.mongoDbRepo.UpdatePartialVoting(ctx, map[string]interface{}{
		"id": voting.ID,
	}, map[string]interface{}{
		"seriesId":         voting.SeriesID,
		"title":            voting.Title,
		"startDate":        voting.StartDate,
		"endDate":          voting.EndDate,
		"performancePoint": voting.PerformancePoint,
		"banner":           voting.Banner,
		"status":           voting.Status,
		"updatedAt":        voting.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, voting.Format())
}

func (u *superadminAppUsecase) DeleteVoting(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check voting
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Voting not found", nil, nil)
	}

	// timestamp
	now := time.Now()

	// validate status & period
	if voting.Status == mongo_model.VotingStatusActive && now.After(voting.StartDate) && now.Before(voting.EndDate) {
		return helpers.NewResponse(http.StatusBadRequest, "Voting cannot be deleted while it is active", nil, nil)
	}

	// delete voting
	voting.DeletedAt = &now

	err = u.mongoDbRepo.UpdatePartialVoting(ctx, map[string]interface{}{
		"id": voting.ID,
	}, map[string]interface{}{
		"deletedAt": voting.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// mark media unused in bg
	go u.markMediaAsUnusedByIds(context.Background(), []string{voting.Banner.ID})

	return helpers.NewResponse(http.StatusOK, "Delete voting success", nil, nil)
}

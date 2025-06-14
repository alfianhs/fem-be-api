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

func (u *superadminAppUsecase) GetSeriesList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if queryParam.Get("status") != "" {
		status := queryParam.Get("status")
		statusInt, err := strconv.Atoi(status)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Invalid status", nil, nil)
		}
		fetchOptions["status"] = mongo_model.SeriesStatus(statusInt)
	}
	if queryParam.Get("seasonId") != "" {
		fetchOptions["seasonId"] = queryParam.Get("seasonId")
	}
	if queryParam.Get("search") != "" {
		fetchOptions["name"] = queryParam.Get("search")
	}

	// count total
	total := u.mongoDbRepo.CountSeries(ctx, fetchOptions)
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
	cur, err := u.mongoDbRepo.FetchListSeries(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var series []mongo_model.Series
	for cur.Next(ctx) {
		row := mongo_model.Series{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Series Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		series = append(series, row)
	}

	// extract season ids
	seasonIds := helpers.ExtractIds(series, func(m mongo_model.Series) string {
		return m.SeasonID
	})

	// fetch season data
	seasons, err := u.mongoDbRepo.FetchListSeason(ctx, map[string]interface{}{
		"ids": seasonIds,
	})
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

	// extract venue ids
	venueIds := helpers.ExtractIds(series, func(m mongo_model.Series) string {
		return m.VenueID
	})

	// fetch venue data
	venues, err := u.mongoDbRepo.FetchListVenue(ctx, map[string]interface{}{
		"ids": venueIds,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer venues.Close(ctx)

	venueMap := make(map[string]mongo_model.Venue)
	for venues.Next(ctx) {
		var venue mongo_model.Venue
		if err := venues.Decode(&venue); err != nil {
			logrus.Error("Venue Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		venueMap[venue.ID.Hex()] = venue
	}

	// format series data
	var list []interface{}
	for _, row := range series {
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

		venue, ok := venueMap[row.VenueID]
		if ok {
			row.Venue = mongo_model.VenueFK{
				ID:   row.VenueID,
				Name: venue.Name,
			}
		} else {
			row.Venue = mongo_model.VenueFK{
				ID:   row.VenueID,
				Name: "",
			}
			logrus.Error("Venue with id ", row.VenueID, " not found")
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

func (u *superadminAppUsecase) GetSeriesDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// fetch data
	row, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if row == nil {
		return helpers.NewResponse(http.StatusNotFound, "Series not found", nil, nil)
	}

	// fetch season data
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": row.SeasonID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusNotFound, "Season not found", nil, nil)
	}

	row.Season = mongo_model.SeasonFK{
		ID:   row.SeasonID,
		Name: season.Name,
	}

	// fetch venue data
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": row.VenueID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusNotFound, "Venue not found", nil, nil)
	}

	row.Venue = mongo_model.VenueFK{
		ID:   row.VenueID,
		Name: venue.Name,
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, row.Format())
}

func (u *superadminAppUsecase) CreateSeries(ctx context.Context, payload request.SeriesCreateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	if payload.VenueID == "" {
		errValidation["venueId"] = "Venue ID field is required"
	}
	if payload.Price <= 0 {
		errValidation["price"] = "Price field is required"
	}
	if payload.StartDate == "" {
		errValidation["startDate"] = "Start date field is required"
	} else {
		_, err := time.Parse(time.RFC3339, payload.StartDate)
		if err != nil {
			errValidation["startDate"] = "Start date format is invalid"
		}
	}
	if payload.EndDate == "" {
		errValidation["endDate"] = "End date field is required"
	} else {
		_, err := time.Parse(time.RFC3339, payload.EndDate)
		if err != nil {
			errValidation["endDate"] = "End date format is invalid"
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation error", errValidation, nil)
	}

	// check active season
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"status": mongo_model.SeasonStatusActive,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Active season not found", nil, nil)
	}

	// check venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": payload.VenueID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
	}

	// prepare series data
	now := time.Now()
	startDate, _ := time.Parse(time.RFC3339, payload.StartDate)
	endDate, _ := time.Parse(time.RFC3339, payload.EndDate)

	// set date to start of day and end of day
	startDate = helpers.SetToStartOfDayUTC(startDate)
	endDate = helpers.SetToEndOfDayUTC(endDate)

	// validate date range
	if startDate.After(endDate) || startDate.Format("2006-01-02") == endDate.Format("2006-01-02") {
		return helpers.NewResponse(http.StatusBadRequest, "Start date must be before end date", nil, nil)
	}

	// create series
	series := mongo_model.Series{
		ID:        primitive.NewObjectID(),
		SeasonID:  season.ID.Hex(),
		Season:    mongo_model.SeasonFK{ID: season.ID.Hex(), Name: season.Name},
		VenueID:   payload.VenueID,
		Venue:     mongo_model.VenueFK{ID: payload.VenueID, Name: venue.Name},
		Name:      payload.Name,
		Price:     payload.Price,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    mongo_model.SeriesStatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save
	err = u.mongoDbRepo.CreateOneSeries(ctx, &series)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Success", nil, series.Format())
}

func (u *superadminAppUsecase) UpdateSeries(ctx context.Context, id string, payload request.SeriesUpdateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check series
	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
	}

	// update if not empty
	if payload.Name != "" {
		series.Name = payload.Name
	}
	if payload.VenueID != "" {
		venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
			"id": payload.VenueID,
		})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		if venue == nil {
			return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
		}
		series.VenueID = payload.VenueID
	}
	if payload.Price > 0 {
		series.Price = payload.Price
	}
	if payload.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, payload.StartDate)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Start date format is invalid", nil, nil)
		}
		// set date to start of day
		startDate = helpers.SetToStartOfDayUTC(startDate)
		series.StartDate = startDate
	}
	if payload.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, payload.EndDate)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "End date format is invalid", nil, nil)
		}
		// set date to end of day
		endDate = helpers.SetToEndOfDayUTC(endDate)
		series.EndDate = endDate
	}
	if payload.Status != nil {
		series.Status = *payload.Status
	}

	// validate date range
	if series.StartDate.After(series.EndDate) || series.StartDate.Format("2006-01-02") == series.EndDate.Format("2006-01-02") {
		return helpers.NewResponse(http.StatusBadRequest, "Start date must be before end date", nil, nil)
	}

	// update series
	series.UpdatedAt = time.Now()
	err = u.mongoDbRepo.UpdatePartialSeries(ctx, map[string]interface{}{
		"id": id,
	}, map[string]interface{}{
		"name":      series.Name,
		"venueId":   series.VenueID,
		"price":     series.Price,
		"startDate": series.StartDate,
		"endDate":   series.EndDate,
		"status":    series.Status,
		"updatedAt": series.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, series.Format())
}

func (u *superadminAppUsecase) DeleteSeries(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check series
	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
	}

	// delete series
	now := time.Now()
	series.DeletedAt = &now

	// delete series
	err = u.mongoDbRepo.UpdatePartialSeries(ctx, map[string]interface{}{
		"id": id,
	}, map[string]interface{}{
		"deletedAt": series.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, nil)
}

func (u *superadminAppUsecase) updateSeriesMatchCount(ctx context.Context, id string) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check tickets with series id
	cur, err := u.mongoDbRepo.FetchListTicket(ctx, map[string]interface{}{
		"seriesId": id,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	defer cur.Close(ctx)

	var tickets []*mongo_model.Ticket
	for cur.Next(ctx) {
		var ticket mongo_model.Ticket
		if err := cur.Decode(&ticket); err != nil {
			logrus.Error("Ticket Decode:", err)
			return
		}
		tickets = append(tickets, &ticket)
	}

	// count matchs
	matchCount := 0
	for _, ticket := range tickets {
		matchCount += len(ticket.Matchs)
	}

	// update series
	err = u.mongoDbRepo.UpdatePartialSeries(ctx, map[string]interface{}{
		"id": id,
	}, map[string]interface{}{
		"matchCount": matchCount,
		"updatedAt":  time.Now(),
	})
	if err != nil {
		logrus.Error(err)
		return
	}
}

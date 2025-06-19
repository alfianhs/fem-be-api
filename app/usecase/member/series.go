package member_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (u *memberAppUsecase) GetSeriesList(ctx context.Context, queryParam url.Values) helpers.Response {
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

func (u *memberAppUsecase) GetSeriesDetail(ctx context.Context, id string) helpers.Response {
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

func (u *memberAppUsecase) GetSeriesListWithTickets(ctx context.Context, queryParam url.Values) helpers.Response {
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

	// extract series ids
	seriesIds := helpers.ExtractIds(series, func(m mongo_model.Series) string {
		return m.ID.Hex()
	})

	// fetch tickets for series
	tickets, err := u.getTicketsBySeriesIDs(ctx, seriesIds)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// group ticket
	groupedTickets := map[string][]mongo_model.Ticket{}
	for _, t := range tickets {
		groupedTickets[t.SeriesID] = append(groupedTickets[t.SeriesID], t)
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

		row.Tickets = groupedTickets[row.ID.Hex()]

		list = append(list, row.Format())
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		Limit: limit,
		Page:  page,
		Total: total,
		List:  list,
	})
}

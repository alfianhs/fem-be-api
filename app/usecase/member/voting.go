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

func (u *memberAppUsecase) GetVotingList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, offset, limit := helpers.GetOffsetLimit(queryParam)
	opts := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	if seriesId := queryParam.Get("seriesId"); seriesId != "" {
		opts["seriesId"] = seriesId
	}
	if s := queryParam.Get("status"); s != "" {
		statusInt, err := strconv.Atoi(s)
		if err != nil {
			return helpers.NewResponse(http.StatusBadRequest, "Invalid status", nil, nil)
		}
		opts["status"] = mongo_model.VotingStatus(statusInt)
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

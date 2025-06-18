package member_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

func (u *memberAppUsecase) GetTicketsList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if queryParam.Get("seriesId") != "" {
		fetchOptions["seriesId"] = queryParam.Get("seriesId")
	}

	// count total
	total := u.mongoDbRepo.CountTicket(ctx, fetchOptions)
	if total == 0 {
		return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
			Limit: limit,
			Page:  page,
			Total: total,
			List:  make([]interface{}, 0),
		})
	}

	// sorting
	if queryParam.Get("sort") != "" {
		fetchOptions["sort"] = queryParam.Get("sort")
	}
	if queryParam.Get("dir") != "" {
		fetchOptions["dir"] = queryParam.Get("dir")
	}

	// fetching
	cur, err := u.mongoDbRepo.FetchListTicket(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var tickets []mongo_model.Ticket
	for cur.Next(ctx) {
		row := mongo_model.Ticket{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Ticket Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		tickets = append(tickets, row)
	}

	// check venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": tickets[0].Matchs[0].VenueID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		venue = &mongo_model.Venue{}
	}

	// check season team without duplicate
	seasonTeamIdSet := make(map[string]struct{})
	for _, ticket := range tickets {
		for _, match := range ticket.Matchs {
			seasonTeamIdSet[match.HomeSeasonTeamID] = struct{}{}
			seasonTeamIdSet[match.AwaySeasonTeamID] = struct{}{}
		}
	}

	// set ids to slice
	seasonTeamIds := make([]string, 0, len(seasonTeamIdSet))
	for id := range seasonTeamIdSet {
		seasonTeamIds = append(seasonTeamIds, id)
	}

	// fetch season team
	seasonTeamCur, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, map[string]interface{}{
		"ids": seasonTeamIds,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer seasonTeamCur.Close(ctx)

	// map season team
	seasonTeamMap := make(map[string]mongo_model.SeasonTeam)
	for seasonTeamCur.Next(ctx) {
		row := mongo_model.SeasonTeam{}
		err := seasonTeamCur.Decode(&row)
		if err != nil {
			logrus.Error("SeasonTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		seasonTeamMap[row.ID.Hex()] = row
	}

	// set list
	var list []interface{}
	for _, ticket := range tickets {
		for i := range ticket.Matchs {
			ticket.Matchs[i].Venue = mongo_model.VenueFK{
				ID:   venue.ID.Hex(),
				Name: venue.Name,
			}

			homeSeasonTeam, ok := seasonTeamMap[ticket.Matchs[i].HomeSeasonTeamID]
			if ok {
				ticket.Matchs[i].HomeSeasonTeam = mongo_model.SeasonTeamFK{
					ID:       homeSeasonTeam.ID.Hex(),
					SeasonID: homeSeasonTeam.SeasonID,
					TeamID:   homeSeasonTeam.Team.ID,
					Team: mongo_model.TeamFK{
						ID:   homeSeasonTeam.Team.ID,
						Name: homeSeasonTeam.Team.Name,
						Logo: homeSeasonTeam.Team.Logo,
					},
				}
			} else {
				ticket.Matchs[i].HomeSeasonTeam = mongo_model.SeasonTeamFK{}
				logrus.Error("Home Season Team " + ticket.Matchs[i].HomeSeasonTeamID + " not found")
			}
			awaySeasonTeam, ok := seasonTeamMap[ticket.Matchs[i].AwaySeasonTeamID]
			if ok {
				ticket.Matchs[i].AwaySeasonTeam = mongo_model.SeasonTeamFK{
					ID:       awaySeasonTeam.ID.Hex(),
					SeasonID: awaySeasonTeam.SeasonID,
					TeamID:   awaySeasonTeam.Team.ID,
					Team: mongo_model.TeamFK{
						ID:   awaySeasonTeam.Team.ID,
						Name: awaySeasonTeam.Team.Name,
						Logo: awaySeasonTeam.Team.Logo,
					},
				}
			} else {
				ticket.Matchs[i].AwaySeasonTeam = mongo_model.SeasonTeamFK{}
				logrus.Error("Away Season Team " + ticket.Matchs[i].AwaySeasonTeamID + " not found")
			}
		}
		list = append(list, ticket.Format())
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		Limit: limit,
		Page:  page,
		Total: total,
		List:  list,
	})
}

func (u *memberAppUsecase) GetTicketDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check ticket
	ticket, err := u.mongoDbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if ticket == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Ticket not found", nil, nil)
	}

	// check venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": ticket.Matchs[0].VenueID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		venue = &mongo_model.Venue{}
	}

	// check season team without duplicate
	seasonTeamIdSet := make(map[string]struct{})
	for _, match := range ticket.Matchs {
		seasonTeamIdSet[match.HomeSeasonTeamID] = struct{}{}
		seasonTeamIdSet[match.AwaySeasonTeamID] = struct{}{}
	}

	// set ids to slice
	seasonTeamIds := make([]string, 0, len(seasonTeamIdSet))
	for id := range seasonTeamIdSet {
		seasonTeamIds = append(seasonTeamIds, id)
	}

	// fetch season team
	seasonTeamCur, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, map[string]interface{}{
		"ids": seasonTeamIds,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer seasonTeamCur.Close(ctx)

	// map season team
	seasonTeamMap := make(map[string]mongo_model.SeasonTeam)
	for seasonTeamCur.Next(ctx) {
		row := mongo_model.SeasonTeam{}
		err := seasonTeamCur.Decode(&row)
		if err != nil {
			logrus.Error("SeasonTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		seasonTeamMap[row.ID.Hex()] = row
	}

	// set detail ticket
	for i := range ticket.Matchs {
		ticket.Matchs[i].Venue = mongo_model.VenueFK{
			ID:   venue.ID.Hex(),
			Name: venue.Name,
		}

		homeSeasonTeam, ok := seasonTeamMap[ticket.Matchs[i].HomeSeasonTeamID]
		if ok {
			ticket.Matchs[i].HomeSeasonTeam = mongo_model.SeasonTeamFK{
				ID:       homeSeasonTeam.ID.Hex(),
				SeasonID: homeSeasonTeam.SeasonID,
				TeamID:   homeSeasonTeam.Team.ID,
				Team: mongo_model.TeamFK{
					ID:   homeSeasonTeam.Team.ID,
					Name: homeSeasonTeam.Team.Name,
					Logo: homeSeasonTeam.Team.Logo,
				},
			}
		} else {
			ticket.Matchs[i].HomeSeasonTeam = mongo_model.SeasonTeamFK{}
			logrus.Error("Home Season Team " + ticket.Matchs[i].HomeSeasonTeamID + " not found")
		}
		awaySeasonTeam, ok := seasonTeamMap[ticket.Matchs[i].AwaySeasonTeamID]
		if ok {
			ticket.Matchs[i].AwaySeasonTeam = mongo_model.SeasonTeamFK{
				ID:       awaySeasonTeam.ID.Hex(),
				SeasonID: awaySeasonTeam.SeasonID,
				TeamID:   awaySeasonTeam.Team.ID,
				Team: mongo_model.TeamFK{
					ID:   awaySeasonTeam.Team.ID,
					Name: awaySeasonTeam.Team.Name,
					Logo: awaySeasonTeam.Team.Logo,
				},
			}
		} else {
			ticket.Matchs[i].AwaySeasonTeam = mongo_model.SeasonTeamFK{}
			logrus.Error("Away Season Team " + ticket.Matchs[i].AwaySeasonTeamID + " not found")
		}
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, ticket.Format())
}

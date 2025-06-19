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

func (u *superadminAppUsecase) GetTicketsList(ctx context.Context, queryParam url.Values) helpers.Response {
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

func (u *superadminAppUsecase) GetTicketDetail(ctx context.Context, id string) helpers.Response {
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

func (u *superadminAppUsecase) CreateOrUpdateTicket(ctx context.Context, payload request.TicketCreateOrUpdateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate payload
	dateSet := make(map[string]struct{})
	errValidation := make(map[string]string)
	if payload.SeriesID == "" {
		errValidation["seriesId"] = "Series ID is required"
	}
	if len(payload.Tickets) == 0 {
		errValidation["tickets"] = "Tickets is required"
	}
	for i, ticket := range payload.Tickets {
		if ticket.Name == "" {
			errValidation["tickets["+strconv.Itoa(i)+"].name"] = "Name is required"
		}
		if ticket.Price <= 0 {
			errValidation["tickets["+strconv.Itoa(i)+"].price"] = "Price is required"
		}
		if ticket.Quota <= 0 {
			errValidation["tickets["+strconv.Itoa(i)+"].quota"] = "Quota is required"
		}
		if ticket.Date == "" {
			errValidation["tickets["+strconv.Itoa(i)+"].date"] = "Date is required"
		} else {
			// date validation
			date, err := time.Parse(time.RFC3339, ticket.Date)
			if err != nil {
				errValidation["tickets["+strconv.Itoa(i)+"].date"] = "Date is invalid"
				continue
			}
			dateStr := date.Format("2006-01-02")
			if _, exists := dateSet[dateStr]; exists {
				errValidation["tickets["+strconv.Itoa(i)+"].date"] = "Duplicate date is not allowed"
			} else {
				dateSet[dateStr] = struct{}{}
			}
		}
		if len(ticket.Matchs) == 0 {
			errValidation["tickets["+strconv.Itoa(i)+"].matchs"] = "Matchs is required"
		}
		timeSet := make(map[string]struct{})
		for j, match := range ticket.Matchs {
			if match.HomeSeasonTeamID == "" {
				errValidation["tickets["+strconv.Itoa(i)+"].matchs["+strconv.Itoa(j)+"].homeSeasonTeamId"] = "Home Season Team ID is required"
			}
			if match.AwaySeasonTeamID == "" {
				errValidation["tickets["+strconv.Itoa(i)+"].matchs["+strconv.Itoa(j)+"].awaySeasonTeamId"] = "Away Season Team ID is required"
			}
			if match.HomeSeasonTeamID == match.AwaySeasonTeamID {
				errValidation["tickets["+strconv.Itoa(i)+"].matchs["+strconv.Itoa(j)+"].homeSeasonTeamId"] = "Home Season Team ID and Away Season Team ID cannot be the same"
			}
			if match.Time == "" {
				errValidation["tickets["+strconv.Itoa(i)+"].matchs["+strconv.Itoa(j)+"].time"] = "Time is required"
			} else {
				// time validation
				_, err := time.Parse("15:04", match.Time)
				if err != nil {
					errValidation["tickets["+strconv.Itoa(i)+"].matchs["+strconv.Itoa(j)+"].time"] = "Time is invalid"
					continue
				}
				if _, exists := timeSet[match.Time]; exists {
					errValidation["tickets["+strconv.Itoa(i)+"].matchs["+strconv.Itoa(j)+"].time"] = "Duplicate match time in this ticket is not allowed"
				} else {
					timeSet[match.Time] = struct{}{}
				}
			}
		}
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check series
	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": payload.SeriesID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
	}

	// check venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": series.VenueID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
	}

	// check season team without duplicate
	seasonTeamIdSet := make(map[string]struct{})
	for _, ticketPayload := range payload.Tickets {
		for _, match := range ticketPayload.Matchs {
			seasonTeamIdSet[match.HomeSeasonTeamID] = struct{}{}
			seasonTeamIdSet[match.AwaySeasonTeamID] = struct{}{}
		}
	}

	// set ids to slice
	seasonTeamIds := make([]string, 0, len(seasonTeamIdSet))
	for id := range seasonTeamIdSet {
		seasonTeamIds = append(seasonTeamIds, id)
	}
	seasonTeams, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, map[string]interface{}{
		"ids": seasonTeamIds,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer seasonTeams.Close(ctx)

	// map season team
	seasonTeamMap := make(map[string]mongo_model.SeasonTeam)
	for seasonTeams.Next(ctx) {
		var seasonTeam mongo_model.SeasonTeam
		if err := seasonTeams.Decode(&seasonTeam); err != nil {
			logrus.Error("SeasonTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		seasonTeamMap[seasonTeam.ID.Hex()] = seasonTeam
	}

	var createdTickets []*mongo_model.Ticket
	var updatedTickets []mongo_model.Ticket
	now := time.Now()
	for _, ticketPayload := range payload.Tickets {
		// set date to start of day
		date, _ := time.Parse(time.RFC3339, ticketPayload.Date)
		date = helpers.SetToStartOfDayWIB(date)

		if date.Before(series.StartDate) || date.After(series.EndDate) {
			return helpers.NewResponse(http.StatusBadRequest, "Date must be between "+series.StartDate.Format("2006-01-02")+
				" and "+series.EndDate.Format("2006-01-02"), nil, nil)
		}

		// set ticket
		ticket := &mongo_model.Ticket{
			Name:     ticketPayload.Name,
			SeriesID: payload.SeriesID,
			Date:     date,
			Price:    ticketPayload.Price,
			Matchs:   []mongo_model.TicketMatch{},
			Quota: mongo_model.TicketQuota{
				Stock: ticketPayload.Quota,
			},
			UpdatedAt: now,
		}

		// set matchs
		for _, matchPayload := range ticketPayload.Matchs {
			_, ok := seasonTeamMap[matchPayload.HomeSeasonTeamID]
			if !ok {
				return helpers.NewResponse(http.StatusBadRequest, "Home Season Team "+matchPayload.HomeSeasonTeamID+" not found", nil, nil)
			}
			_, ok = seasonTeamMap[matchPayload.AwaySeasonTeamID]
			if !ok {
				return helpers.NewResponse(http.StatusBadRequest, "Away Season Team "+matchPayload.AwaySeasonTeamID+" not found", nil, nil)
			}

			ticket.Matchs = append(ticket.Matchs, mongo_model.TicketMatch{
				HomeSeasonTeamID: matchPayload.HomeSeasonTeamID,
				AwaySeasonTeamID: matchPayload.AwaySeasonTeamID,
				VenueID:          venue.ID.Hex(),
				Time:             matchPayload.Time,
			})
		}

		// if id is empty, create ticket
		if ticketPayload.ID == "" {
			ticket.ID = primitive.NewObjectID()
			ticket.CreatedAt = now
			createdTickets = append(createdTickets, ticket)
		} else {
			// if id is not empty, update ticket
			existingTicket, err := u.mongoDbRepo.FetchOneTicket(ctx, map[string]interface{}{
				"id": ticketPayload.ID,
			})
			if err != nil {
				return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
			}
			if existingTicket == nil {
				return helpers.NewResponse(http.StatusBadRequest, "Ticket "+ticketPayload.ID+" not found", nil, nil)
			}

			// save updated ticket
			err = u.mongoDbRepo.UpdatePartialTicket(ctx, map[string]interface{}{
				"id": ticketPayload.ID,
			}, map[string]interface{}{
				"name":      ticket.Name,
				"date":      ticket.Date,
				"price":     ticket.Price,
				"quota":     ticket.Quota,
				"matchs":    ticket.Matchs,
				"updatedAt": ticket.UpdatedAt,
			})
			if err != nil {
				return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
			}
			updatedTickets = append(updatedTickets, *existingTicket)
		}
	}

	// create many ticket if createdTickets is not empty
	if len(createdTickets) > 0 {
		err = u.mongoDbRepo.CreateManyTicket(ctx, createdTickets)
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
	}

	// update match count in related series in bg
	go u.updateSeriesMatchCount(context.Background(), series.ID.Hex())

	// return response
	return helpers.NewResponse(http.StatusOK, "Success", nil, map[string]interface{}{
		"created": createdTickets,
		"updated": updatedTickets,
	})
}

func (u *superadminAppUsecase) DeleteTicket(ctx context.Context, id string) helpers.Response {
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

	// delete ticket
	now := time.Now()
	ticket.DeletedAt = &now
	err = u.mongoDbRepo.UpdatePartialTicket(ctx, map[string]interface{}{
		"id": ticket.ID,
	}, map[string]interface{}{
		"deletedAt": ticket.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// update match count in related series in bg
	go u.updateSeriesMatchCount(context.Background(), ticket.SeriesID)

	return helpers.NewResponse(http.StatusOK, "Success", nil, nil)
}

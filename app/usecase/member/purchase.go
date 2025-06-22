package member_usecase

import (
	mongo_model "app/domain/model/mongo"
	xendit_model "app/domain/model/xendit"
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *memberAppUsecase) CreatePurchase(ctx context.Context, claim jwt_helpers.MemberJWTClaims, payload request.CreatePurchaseRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.ProductId == "" {
		errValidation["productId"] = "Product ID field is required"
	}
	if payload.Amount <= 0 {
		errValidation["amount"] = "Amount field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// max amount 4
	if payload.Amount > 4 {
		return helpers.NewResponse(http.StatusBadRequest, "Max amount buy is 4", nil, nil)
	}

	// check ticket
	ticket, err := u.mongoDbRepo.FetchOneTicket(ctx, map[string]interface{}{
		"id": payload.ProductId,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if ticket == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Ticket not found", nil, nil)
	}

	ticket.Format()

	// check quota
	if ticket.Quota.Remaining < payload.Amount {
		return helpers.NewResponse(http.StatusBadRequest, "Ticket quota is not enough", nil, nil)
	}

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	if member.Phone == nil {
		phone := ""
		member.Phone = &phone
	}

	// check series
	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": ticket.SeriesID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
	}

	// check season
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": series.SeasonID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season not found", nil, nil)
	}

	// create new purchase
	pricePcs := ticket.Price
	grandTotal := pricePcs * float64(payload.Amount)
	now := time.Now()

	// count purchase today for external ID
	purchaseCount := u.mongoDbRepo.CountPurchase(ctx, map[string]interface{}{
		"memberId": member.ID.Hex(),
		"today":    true,
	})

	// generate external ID
	externalId := helpers.GenerateInvoiceExternalId(purchaseCount)

	newPurchase := mongo_model.Purchase{
		ID: primitive.NewObjectID(),
		Member: mongo_model.MemberPurchaseFK{
			ID:    member.ID.Hex(),
			Name:  member.Name,
			Email: member.Email,
			Phone: *member.Phone,
		},
		SeasonId: season.ID.Hex(),
		Season: mongo_model.SeasonFK{
			ID:   season.ID.Hex(),
			Name: season.Name,
		},
		SeriesId: series.ID.Hex(),
		Series: mongo_model.SeriesFK{
			ID:   series.ID.Hex(),
			Name: series.Name,
		},
		Tickets: []mongo_model.TicketFK{
			{
				ID:      payload.ProductId,
				Name:    ticket.Name,
				Date:    ticket.Date,
				VenueID: series.VenueID,
			},
		},
		Invoice: mongo_model.Invoice{
			InvoiceExternalID: externalId,
		},
		Amount:            payload.Amount,
		Price:             pricePcs,
		GrandTotal:        grandTotal,
		IsCheckoutPackage: false,
		Status:            mongo_model.PurchaseStatusPending,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if pricePcs > 0 {
		// generate snap link
		result, err := u.xenditRepo.GenereteSnapLink(ctx, newPurchase)
		if err != nil || result.Status != http.StatusOK {
			if result.Status != 0 {
				return helpers.NewResponse(http.StatusBadRequest, result.Message, nil, nil)
			}

			return helpers.NewResponse(http.StatusBadRequest, err.Error(), nil, nil)
		}
		respDataXendit, _ := result.Data.(xendit_model.XenditSnapLinkSuccessResponse)

		newPurchase.Invoice.InvoiceID = respDataXendit.ID
		newPurchase.Invoice.InvoiceUrl = respDataXendit.InvoiceURL
		newPurchase.Invoice.MerchantName = respDataXendit.MerchantName
		newPurchase.ExpiresAt = respDataXendit.ExpiryDate
	}

	// save purchase
	err = u.mongoDbRepo.CreateOnePurchase(ctx, &newPurchase)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// update ticket quota bg
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
		defer cancel()

		for _, p := range newPurchase.Tickets {
			err := u.mongoDbRepo.IncrementOneTicket(ctx, p.ID, map[string]int64{
				"quota.used": +newPurchase.Amount,
			})
			if err != nil {
				logrus.Error("IncrementOneTicket:", err)
			}
		}
	}()

	return helpers.NewResponse(http.StatusOK, "Purchase success", nil, newPurchase)
}

func (u *memberAppUsecase) CreatePackagePurchase(ctx context.Context, claim jwt_helpers.MemberJWTClaims, payload request.CreatePurchaseRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.ProductId == "" {
		errValidation["productId"] = "Product ID field is required"
	}
	if payload.Amount <= 0 {
		errValidation["amount"] = "Amount field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// max amount 4
	if payload.Amount > 4 {
		return helpers.NewResponse(http.StatusBadRequest, "Max amount buy is 4", nil, nil)
	}

	// check series
	series, err := u.mongoDbRepo.FetchOneSeries(ctx, map[string]interface{}{
		"id": payload.ProductId,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if series == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Series not found", nil, nil)
	}

	// check ticket
	ticketCur, err := u.mongoDbRepo.FetchListTicket(ctx, map[string]interface{}{
		"seriesId": payload.ProductId,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer ticketCur.Close(ctx)

	tickets := []mongo_model.Ticket{}
	for ticketCur.Next(ctx) {
		row := mongo_model.Ticket{}
		err := ticketCur.Decode(&row)
		if err != nil {
			logrus.Error("Ticket Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		tickets = append(tickets, row)
	}

	if len(tickets) == 0 {
		return helpers.NewResponse(http.StatusBadRequest, "Ticket not found", nil, nil)
	}

	// check quota
	var ticketsFK []mongo_model.TicketFK
	for _, ticket := range tickets {
		ticket.Format()
		if ticket.Quota.Remaining < payload.Amount {
			return helpers.NewResponse(http.StatusBadRequest, "Ticket quota is not enough", nil, nil)
		}
		ticketsFK = append(ticketsFK, mongo_model.TicketFK{
			ID:      ticket.ID.Hex(),
			Name:    ticket.Name,
			Date:    ticket.Date,
			VenueID: series.VenueID,
		})
	}

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	if member.Phone == nil {
		phone := ""
		member.Phone = &phone
	}

	// check season
	season, err := u.mongoDbRepo.FetchOneSeason(ctx, map[string]interface{}{
		"id": series.SeasonID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if season == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Season not found", nil, nil)
	}

	// create new purchase
	pricePcs := series.Price
	grandTotal := pricePcs * float64(payload.Amount)
	now := time.Now()

	// count purchase today for external ID
	purchaseCount := u.mongoDbRepo.CountPurchase(ctx, map[string]interface{}{
		"memberId": member.ID.Hex(),
		"today":    true,
	})

	// generate external ID
	externalId := helpers.GenerateInvoiceExternalId(purchaseCount)

	newPurchase := mongo_model.Purchase{
		ID: primitive.NewObjectID(),
		Member: mongo_model.MemberPurchaseFK{
			ID:    member.ID.Hex(),
			Name:  member.Name,
			Email: member.Email,
			Phone: *member.Phone,
		},
		SeasonId: season.ID.Hex(),
		Season: mongo_model.SeasonFK{
			ID:   season.ID.Hex(),
			Name: season.Name,
		},
		SeriesId: series.ID.Hex(),
		Series: mongo_model.SeriesFK{
			ID:   series.ID.Hex(),
			Name: series.Name,
		},
		Tickets: ticketsFK,
		Invoice: mongo_model.Invoice{
			InvoiceExternalID: externalId,
		},
		Amount:            payload.Amount,
		Price:             pricePcs,
		GrandTotal:        grandTotal,
		IsCheckoutPackage: true,
		Status:            mongo_model.PurchaseStatusPending,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if pricePcs > 0 {
		// generate snap link
		result, err := u.xenditRepo.GenereteSnapLink(ctx, newPurchase)
		if err != nil || result.Status != http.StatusOK {
			if result.Status != 0 {
				return helpers.NewResponse(http.StatusBadRequest, result.Message, nil, nil)
			}

			return helpers.NewResponse(http.StatusBadRequest, err.Error(), nil, nil)
		}
		respDataXendit, _ := result.Data.(xendit_model.XenditSnapLinkSuccessResponse)

		newPurchase.Invoice.InvoiceID = respDataXendit.ID
		newPurchase.Invoice.InvoiceUrl = respDataXendit.InvoiceURL
		newPurchase.Invoice.MerchantName = respDataXendit.MerchantName
		newPurchase.ExpiresAt = respDataXendit.ExpiryDate
	}

	// save purchase
	err = u.mongoDbRepo.CreateOnePurchase(ctx, &newPurchase)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// update ticket quota bg
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
		defer cancel()

		for _, p := range newPurchase.Tickets {
			err := u.mongoDbRepo.IncrementOneTicket(ctx, p.ID, map[string]int64{
				"quota.used": +newPurchase.Amount,
			})
			if err != nil {
				logrus.Error("IncrementOneTicket:", err)
			}
		}
	}()

	return helpers.NewResponse(http.StatusOK, "Purchase success", nil, newPurchase)
}

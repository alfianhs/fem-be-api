package admin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (u *adminAppUsecase) GetListTicketPurchasesIsUsedToday(ctx context.Context) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	fetchOptions := map[string]interface{}{
		"today":  true,
		"isUsed": true,
	}

	// count
	total := u.mongoDbRepo.CountTicketPurchase(ctx, fetchOptions)
	if total == 0 {
		return helpers.NewResponse(http.StatusOK, "Success", nil, map[string]interface{}{
			"list": []interface{}{},
			"total": 0,
		})
	}

	// fetch ticket purchases
	cur, err := u.mongoDbRepo.FetchListTicketPurchase(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var list []interface{}
	for cur.Next(ctx) {
		row := mongo_model.TicketPurchase{}
		err = cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListPurhcase Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		list = append(list, row)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

func (u *adminAppUsecase) ScanTicketPurchase(ctx context.Context, payload request.ScanTicketPurchaseRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Code == "" {
		errValidation["code"] = "Code field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check ticket
	ticketPurchase, err := u.mongoDbRepo.FetchOneTicketPurchase(ctx, map[string]interface{}{
		"code": payload.Code,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if ticketPurchase == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Ticket not found", nil, nil)
	}

	// validate date
	now := time.Now()
	if !helpers.IsSameDateWIB(ticketPurchase.Ticket.Date, now) {
		return helpers.NewResponse(http.StatusBadRequest, "You can scan this ticket at "+helpers.FormatDateWIB(ticketPurchase.Ticket.Date, "02 January 2006"), nil, nil)
	}

	// update ticket purchase
	ticketPurchase.IsUsed = true
	ticketPurchase.UsedAt = &now
	ticketPurchase.UpdatedAt = now

	err = u.mongoDbRepo.UpdatePartialTicketPurchase(ctx, map[string]interface{}{
		"id": ticketPurchase.ID,
	}, map[string]interface{}{
		"isUsed":    ticketPurchase.IsUsed,
		"usedAt":    ticketPurchase.UsedAt,
		"updatedAt": ticketPurchase.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, ticketPurchase)
}

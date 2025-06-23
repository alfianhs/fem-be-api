package admin_usecase

import (
	"app/domain/request"
	"app/helpers"
	"context"
	"net/http"
	"time"
)

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

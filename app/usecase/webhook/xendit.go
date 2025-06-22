package webhook_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	mailing_helpers "app/helpers/mailing"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *webhookAppUsecase) HandleXenditWebhook(ctx context.Context, payload request.SnapWebhookRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get purchase by external ID and invoice ID
	purchase, err := u.mongoDbRepo.FetchOnePurchase(ctx, map[string]interface{}{
		"invoiceId":  payload.ID,
		"externalId": payload.ExternalID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if purchase == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Purchase not found", nil, nil)
	}

	// check if purchase is already paid
	if purchase.Status == mongo_model.PurchaseStatusPaid {
		return helpers.NewResponse(http.StatusBadRequest, "Purchase already paid", nil, nil)
	}

	// handle the webhook based on the status
	if payload.Status == "PAID" {
		return u.paidPurchase(ctx, payload, purchase)
	} else {
		return u.restoreQuota(ctx, purchase)
	}
}

func (u *webhookAppUsecase) paidPurchase(ctx context.Context, payload request.SnapWebhookRequest, purchase *mongo_model.Purchase) helpers.Response {
	// update purchase
	now := time.Now()
	purchase.Status = mongo_model.PurchaseStatusPaid
	purchase.PaidAt = &payload.PaidAt
	purchase.Invoice = mongo_model.Invoice{
		InvoiceID:          payload.ID,
		InvoiceExternalID:  payload.ExternalID,
		PaymentMethod:      payload.PaymentMethod,
		BankCode:           payload.BankCode,
		PaymentChannel:     payload.PaymentChannel,
		PaymentDestination: payload.PaymentDestination,
		MerchantName:       payload.MerchantName,
	}
	purchase.UpdatedAt = now

	err := u.mongoDbRepo.UpdatePartialPurchase(ctx, map[string]interface{}{
		"id": purchase.ID,
	}, map[string]interface{}{
		"status":    purchase.Status,
		"paidAt":    purchase.PaidAt,
		"invoice":   purchase.Invoice,
		"updatedAt": now,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// get venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": purchase.Tickets[0].VenueID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
	}

	// create array of ticket purchases
	var ticketPurchases []*mongo_model.TicketPurchase
	for _, ticket := range purchase.Tickets {
		for a := 0; a < int(purchase.Amount); a++ {
			ticketPurchase := &mongo_model.TicketPurchase{
				ID:     primitive.NewObjectID(),
				Member: purchase.Member,
				Ticket: ticket,
				Venue: mongo_model.VenueFK{
					ID:   venue.ID.Hex(),
					Name: venue.Name,
				},
				PurchaseID: purchase.ID.Hex(),
				Code:       uuid.NewString(),
				IsUsed:     false,
				UsedAt:     nil,
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			ticketPurchases = append(ticketPurchases, ticketPurchase)
		}
	}

	// create ticket purchases
	err = u.mongoDbRepo.CreateManyTicketPurchase(ctx, ticketPurchases)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// send email to member
	go mailing_helpers.SendTicketPurchase(ticketPurchases)

	return helpers.NewResponse(http.StatusOK, "Ticket purchase generated successfully", nil, purchase)
}

func (u *webhookAppUsecase) restoreQuota(ctx context.Context, purchase *mongo_model.Purchase) helpers.Response {
	// update purchase
	now := time.Now()
	purchase.Status = mongo_model.PurchaseStatusFailed
	purchase.UpdatedAt = now

	err := u.mongoDbRepo.UpdatePartialPurchase(ctx, map[string]interface{}{
		"id": purchase.ID,
	}, map[string]interface{}{
		"status":    purchase.Status,
		"updatedAt": now,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// restore quota for each ticket in bg
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
		defer cancel()

		for _, p := range purchase.Tickets {
			err := u.mongoDbRepo.IncrementOneTicket(ctx, p.ID, map[string]int64{
				"quota.used": purchase.Amount * -1,
			})
			if err != nil {
				logrus.Error("IncrementOneTicket:", err)
			}
		}
	}()

	return helpers.NewResponse(http.StatusOK, "Quota restored successfully", nil, nil)
}

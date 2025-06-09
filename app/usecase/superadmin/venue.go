package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *superadminAppUsecase) GetVenueList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get limit offset
	page, offset, limit := helpers.GetOffsetLimit(queryParam)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if queryParam.Get("search") != "" {
		fetchOptions["name"] = queryParam.Get("search")
	}

	// count total
	total := u.mongoDbRepo.CountVenue(ctx, fetchOptions)
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
	cur, err := u.mongoDbRepo.FetchListVenue(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var list []interface{}
	for cur.Next(ctx) {
		row := mongo_model.Venue{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListVenue Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}

		list = append(list, row)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		Limit: limit,
		Page:  page,
		Total: total,
		List:  list,
	})
}

func (u *superadminAppUsecase) GetVenueDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, venue)
}

func (u *superadminAppUsecase) CreateVenue(ctx context.Context, payload request.VenueCreateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// create venue
	now := time.Now()
	venue := mongo_model.Venue{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// save
	err := u.mongoDbRepo.CreateOneVenue(ctx, &venue)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Create venue success", nil, venue)
}

func (u *superadminAppUsecase) UpdateVenue(ctx context.Context, id string, payload request.VenueUpdateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// validate payload
	errValidation := make(map[string]string)
	if payload.Name == "" {
		errValidation["name"] = "Name field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// get venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
	}

	// update venue
	venue.Name = payload.Name
	venue.UpdatedAt = time.Now()

	// save
	err = u.mongoDbRepo.UpdatePartialVenue(ctx, map[string]interface{}{
		"id": venue.ID,
	}, map[string]interface{}{
		"name":      venue.Name,
		"updatedAt": venue.UpdatedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Update venue success", nil, venue)
}

func (u *superadminAppUsecase) DeleteVenue(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get venue
	venue, err := u.mongoDbRepo.FetchOneVenue(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if venue == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Venue not found", nil, nil)
	}

	// delete venue
	now := time.Now()
	venue.DeletedAt = &now

	// save
	err = u.mongoDbRepo.UpdatePartialVenue(ctx, map[string]interface{}{
		"id": venue.ID,
	}, map[string]interface{}{
		"deletedAt": venue.DeletedAt,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Delete venue success", nil, nil)
}

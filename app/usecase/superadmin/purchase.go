package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

func (u *superadminAppUsecase) GetPurchasesList(ctx context.Context, queryParam url.Values) helpers.Response {
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
		fetchOptions["search"] = queryParam.Get("search")
	}

	// count total
	total := u.mongoDbRepo.CountPurchase(ctx, fetchOptions)
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

	// fetch list
	cur, err := u.mongoDbRepo.FetchListPurchase(ctx, fetchOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var list []interface{}
	for cur.Next(ctx) {
		row := mongo_model.Purchase{}
		err = cur.Decode(&row)
		if err != nil {
			logrus.Error("GetListPurhcase Decode:", err)
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

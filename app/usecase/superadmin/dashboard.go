package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

func (u *superadminAppUsecase) GetDashboard(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// filtering
	purchaseOptions := map[string]interface{}{
		"status": mongo_model.PurchaseStatusPaid,
	}
	seriesOptions := map[string]interface{}{}
	if queryParam.Get("seasonId") != "" {
		purchaseOptions["seasonId"] = queryParam.Get("seasonId")
		seriesOptions["seasonId"] = queryParam.Get("seasonId")
	}

	// sum total match in series
	totalMatch, err := u.mongoDbRepo.SumSeriesMatchCount(ctx, seriesOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// fetch list purchase
	cur, err := u.mongoDbRepo.FetchListPurchase(ctx, purchaseOptions)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, "Error", nil, nil)
	}
	defer cur.Close(ctx)

	var purchases []mongo_model.Purchase
	for cur.Next(ctx) {
		row := mongo_model.Purchase{}
		if err := cur.Decode(&row); err != nil {
			logrus.Error("GetListPurhcase Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, "Error", nil, nil)
		}
		purchases = append(purchases, row)
	}

	totalPurchase := len(purchases)
	totalSeriesPurchase := 0
	totalDayPurchase := 0
	totalIncome := float64(0)
	incomeByMonth := make(map[string]float64)

	// load time location for used timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	for _, purchase := range purchases {
		// sum total income
		totalIncome += purchase.GrandTotal

		// sum total income by month
		month := purchase.CreatedAt.In(loc).Format("Jan")
		if _, ok := incomeByMonth[month]; !ok {
			incomeByMonth[month] = 0
		}
		incomeByMonth[month] += purchase.GrandTotal

		if purchase.IsCheckoutPackage {
			totalSeriesPurchase += 1
		} else {
			totalDayPurchase += 1
		}
	}

	// get all month
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	// final data
	var chartData []map[string]interface{}

	for _, month := range months {
		chartData = append(chartData, map[string]interface{}{
			"month": month,
			"value": incomeByMonth[month],
		})
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, map[string]interface{}{
		"totalMatch":          totalMatch,
		"totalPurchase":       totalPurchase,
		"totalSeriesPurchase": totalSeriesPurchase,
		"totalDayPurchase":    totalDayPurchase,
		"totalIncome":         totalIncome,
		"chartData":           chartData,
	})
}

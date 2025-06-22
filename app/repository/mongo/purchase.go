package mongo_repository

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterPurchase(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if status, ok := options["status"].(mongo_model.PurchaseStatus); ok {
		query["status"] = status
	}
	if memberId, ok := options["memberId"].(string); ok {
		query["member.id"] = memberId
	}
	if today, ok := options["today"].(bool); ok {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		if today {
			query["createdAt"] = bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			}
		}
	}
	if invoiceId, ok := options["invoiceId"].(string); ok {
		query["invoice.invoiceId"] = invoiceId
	}
	if externalId, ok := options["externalId"].(string); ok {
		query["invoice.invoiceExternalId"] = externalId
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListPurchase(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterPurchase(options, true)

	cur, err = r.Conn.Collection(r.purchaseCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListPurchase Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountPurchase(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterPurchase(options, true)

	total, err := r.Conn.Collection(r.purchaseCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountPurchase CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOnePurchase(ctx context.Context, options map[string]interface{}) (row *mongo_model.Purchase, err error) {
	query, _ := generateQueryFilterPurchase(options, false)

	err = r.Conn.Collection(r.purchaseCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOnePurchase FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOnePurchase(ctx context.Context, purchase *mongo_model.Purchase) (err error) {
	_, err = r.Conn.Collection(r.purchaseCollection).InsertOne(ctx, purchase)
	if err != nil {
		logrus.Error("CreatePurchase InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialPurchase(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterPurchase(options, false)

	_, err = r.Conn.Collection(r.purchaseCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialPurchase UpdateOne:", err)
		return
	}

	return
}

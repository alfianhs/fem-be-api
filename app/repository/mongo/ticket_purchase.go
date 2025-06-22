package mongo_repository

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterTicketPurchase(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListTicketPurchase(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterTicketPurchase(options, true)

	cur, err = r.Conn.Collection(r.ticketPurchaseCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListTicketPurchase Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountTicketPurchase(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTicketPurchase(options, true)

	total, err := r.Conn.Collection(r.ticketPurchaseCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicketPurchase CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneTicketPurchase(ctx context.Context, options map[string]interface{}) (row *mongo_model.TicketPurchase, err error) {
	query, _ := generateQueryFilterTicketPurchase(options, false)

	err = r.Conn.Collection(r.ticketPurchaseCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTicketPurchase FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneTicketPurchase(ctx context.Context, ticketPurchase *mongo_model.TicketPurchase) (err error) {
	_, err = r.Conn.Collection(r.ticketPurchaseCollection).InsertOne(ctx, ticketPurchase)
	if err != nil {
		logrus.Error("CreateOneTicketPurchase InsertOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateManyTicketPurchase(ctx context.Context, ticketPurchases []*mongo_model.TicketPurchase) (err error) {
	docs := make([]interface{}, len(ticketPurchases))
	for i, row := range ticketPurchases {
		docs[i] = row
	}

	_, err = r.Conn.Collection(r.ticketPurchaseCollection).InsertMany(ctx, docs)
	if err != nil {
		logrus.Error("CreateManyTicketPurchase InsertMany:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialTicketPurchase(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterTicketPurchase(options, false)

	_, err = r.Conn.Collection(r.ticketPurchaseCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialTicketPurchase UpdateOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) UpdateManyTicketPurchasePartial(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterTicketPurchase(options, false)

	_, err = r.Conn.Collection(r.ticketPurchaseCollection).UpdateMany(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdateManyTicketPurchasePartial UpdateMany:", err)
		return
	}

	return
}

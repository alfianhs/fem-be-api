package mongo_repository

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func generateQueryFilterTicket(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if seriesId, ok := options["seriesId"]; ok {
		query["seriesId"] = seriesId
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListTicket(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterTicket(options, true)

	cur, err = r.Conn.Collection(r.ticketCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListTicket Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountTicket(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTicket(options, true)

	total, err := r.Conn.Collection(r.ticketCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTicket CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneTicket(ctx context.Context, options map[string]interface{}) (row *mongo_model.Ticket, err error) {
	query, _ := generateQueryFilterTicket(options, false)

	err = r.Conn.Collection(r.ticketCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTicket FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneTicket(ctx context.Context, ticket *mongo_model.Ticket) (err error) {
	_, err = r.Conn.Collection(r.ticketCollection).InsertOne(ctx, ticket)
	if err != nil {
		logrus.Error("CreateTicket InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) CreateManyTicket(ctx context.Context, tickets []*mongo_model.Ticket) (err error) {
	docs := make([]interface{}, len(tickets))
	for i, row := range tickets {
		docs[i] = row
	}

	_, err = r.Conn.Collection(r.ticketCollection).InsertMany(ctx, docs)
	if err != nil {
		logrus.Error("CreateManyTicket InsertMany:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialTicket(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterTicket(options, false)

	_, err = r.Conn.Collection(r.ticketCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialTicket UpdateOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) IncrementOneTicket(ctx context.Context, id string, payload map[string]int64) (err error) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Error("Invalid ticket ID:", err)
		return err
	}

	_, err = r.Conn.Collection(r.ticketCollection).UpdateOne(ctx, map[string]any{
		"_id": obj,
	}, bson.M{
		"$inc": payload,
	})
	if err != nil {
		logrus.Error("IncrementOneTicket UpdateOne:", err)
		return
	}
	return
}

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

func generateQueryFilterVenue(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListVenue(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterVenue(options, true)

	cur, err = r.Conn.Collection(r.venueCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListVenue Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountVenue(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterVenue(options, true)

	total, err := r.Conn.Collection(r.venueCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountVenue CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneVenue(ctx context.Context, options map[string]interface{}) (row *mongo_model.Venue, err error) {
	query, _ := generateQueryFilterVenue(options, false)

	err = r.Conn.Collection(r.venueCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneVenue FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneVenue(ctx context.Context, venue *mongo_model.Venue) (err error) {
	_, err = r.Conn.Collection(r.venueCollection).InsertOne(ctx, venue)
	if err != nil {
		logrus.Error("CreateVenue InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialVenue(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterVenue(options, false)

	_, err = r.Conn.Collection(r.venueCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialVenue UpdateOne:", err)
		return
	}

	return
}

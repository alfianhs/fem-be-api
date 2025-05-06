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

func generateQueryFilterPlayer(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListPlayer(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterPlayer(options, true)

	cur, err = r.Conn.Collection(r.playerCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListPlayer Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountPlayer(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterPlayer(options, true)

	total, err := r.Conn.Collection(r.playerCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountPlayer CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOnePlayer(ctx context.Context, options map[string]interface{}) (row *mongo_model.Player, err error) {
	query, _ := generateQueryFilterPlayer(options, false)

	err = r.Conn.Collection(r.playerCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOnePlayer FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOnePlayer(ctx context.Context, player *mongo_model.Player) (err error) {
	_, err = r.Conn.Collection(r.playerCollection).InsertOne(ctx, player)
	if err != nil {
		logrus.Error("CreatePlayer InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialPlayer(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterPlayer(options, false)

	_, err = r.Conn.Collection(r.playerCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialPlayer UpdateOne:", err)
		return
	}

	return
}

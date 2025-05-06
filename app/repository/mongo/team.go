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

func generateQueryFilterTeam(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListTeam(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterTeam(options, true)

	cur, err = r.Conn.Collection(r.teamCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListTeam Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountTeam(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterTeam(options, true)

	total, err := r.Conn.Collection(r.teamCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountTeam CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneTeam(ctx context.Context, options map[string]interface{}) (row *mongo_model.Team, err error) {
	query, _ := generateQueryFilterTeam(options, false)

	err = r.Conn.Collection(r.teamCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneTeam FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneTeam(ctx context.Context, team *mongo_model.Team) (err error) {
	_, err = r.Conn.Collection(r.teamCollection).InsertOne(ctx, team)
	if err != nil {
		logrus.Error("CreateTeam InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialTeam(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterTeam(options, false)

	_, err = r.Conn.Collection(r.teamCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialTeam UpdateOne:", err)
		return
	}

	return
}

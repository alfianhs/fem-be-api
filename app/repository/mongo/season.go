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

func generateQueryFilterSeason(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if status, ok := options["status"].(mongo_model.SeasonStatus); ok {
		query["status"] = status
	}
	if name, ok := options["name"].(string); ok {
		regex := bson.M{
			"$regex": primitive.Regex{
				Pattern: name,
				Options: "i",
			},
		}
		query["name"] = regex
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListSeason(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterSeason(options, true)

	cur, err = r.Conn.Collection(r.seasonCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListSeason Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountSeason(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterSeason(options, true)

	total, err := r.Conn.Collection(r.seasonCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountSeason CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneSeason(ctx context.Context, options map[string]interface{}) (row *mongo_model.Season, err error) {
	query, _ := generateQueryFilterSeason(options, false)

	err = r.Conn.Collection(r.seasonCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneSeason FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneSeason(ctx context.Context, season *mongo_model.Season) (err error) {
	_, err = r.Conn.Collection(r.seasonCollection).InsertOne(ctx, season)
	if err != nil {
		logrus.Error("CreateSeason InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialSeason(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterSeason(options, false)

	_, err = r.Conn.Collection(r.seasonCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialSeason UpdateOne:", err)
		return
	}

	return
}

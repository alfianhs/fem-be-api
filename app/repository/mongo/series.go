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

func generateQueryFilterSeries(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if status, ok := options["status"].(mongo_model.SeriesStatus); ok {
		query["status"] = status
	}
	if seasonId, ok := options["seasonId"].(string); ok {
		query["seasonId"] = seasonId
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

func (r *mongoDbRepo) FetchListSeries(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterSeries(options, true)

	cur, err = r.Conn.Collection(r.seriesCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListSeries Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountSeries(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterSeries(options, true)

	total, err := r.Conn.Collection(r.seriesCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountSeries CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneSeries(ctx context.Context, options map[string]interface{}) (row *mongo_model.Series, err error) {
	query, _ := generateQueryFilterSeries(options, false)

	err = r.Conn.Collection(r.seriesCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneSeries FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneSeries(ctx context.Context, series *mongo_model.Series) (err error) {
	_, err = r.Conn.Collection(r.seriesCollection).InsertOne(ctx, series)
	if err != nil {
		logrus.Error("CreateSeries InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialSeries(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterSeries(options, false)

	_, err = r.Conn.Collection(r.seriesCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialSeries UpdateOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) SumSeriesMatchCount(ctx context.Context, options map[string]interface{}) (total int64, err error) {
	// filter
	query, _ := generateQueryFilterSeries(options, false)

	// pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: query}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalMatchCount", Value: bson.D{{Key: "$sum", Value: "$matchCount"}}},
		}}},
	}

	// aggregate
	cur, err := r.Conn.Collection(r.seriesCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)

	// decode
	if cur.Next(ctx) {
		var result struct {
			TotalMatchCount int64 `bson:"totalMatchCount"`
		}
		if err := cur.Decode(&result); err != nil {
			return 0, err
		}
		total = result.TotalMatchCount
	}

	return total, nil
}

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

func generateQueryFilterSeasonTeam(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if seasonId, ok := options["seasonId"].(string); ok {
		query["seasonId"] = seasonId
	}
	if teamId, ok := options["team.id"].(string); ok {
		query["team.id"] = teamId
	}
	if teamIds, ok := options["team.ids"].([]string); ok {
		query["team.id"] = bson.M{
			"$in": teamIds,
		}
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListSeasonTeam(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterSeasonTeam(options, true)

	cur, err = r.Conn.Collection(r.seasonTeamCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListSeasonTeam Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountSeasonTeam(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterSeasonTeam(options, true)

	total, err := r.Conn.Collection(r.seasonTeamCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountSeasonTeam CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneSeasonTeam(ctx context.Context, options map[string]interface{}) (row *mongo_model.SeasonTeam, err error) {
	query, _ := generateQueryFilterSeasonTeam(options, false)

	err = r.Conn.Collection(r.seasonTeamCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneSeasonTeam FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateManySeasonTeam(ctx context.Context, seasonTeams []*mongo_model.SeasonTeam) (err error) {
	docs := make([]interface{}, len(seasonTeams))
	for i, row := range seasonTeams {
		docs[i] = row
	}

	_, err = r.Conn.Collection(r.seasonTeamCollection).InsertMany(ctx, docs)
	if err != nil {
		logrus.Error("CreateManySeasonTeam InsertMany:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialSeasonTeam(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterSeasonTeam(options, false)

	_, err = r.Conn.Collection(r.seasonTeamCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialSeasonTeam UpdateOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) UpdateManySeasonTeamPartial(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterSeasonTeam(options, false)

	_, err = r.Conn.Collection(r.seasonTeamCollection).UpdateMany(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdateManySeasonTeamPartial UpdateMany:", err)
		return
	}

	return
}

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

func generateQueryFilterSeasonTeamPlayer(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if seasonId, ok := options["seasonId"].(string); ok {
		query["seasonTeam.seasonId"] = seasonId
	}
	if seasonTeamId, ok := options["seasonTeamId"].(string); ok {
		query["seasonTeam.id"] = seasonTeamId
	}
	if teamId, ok := options["teamId"].(string); ok {
		query["seasonTeam.teamId"] = teamId
	}
	if playerId, ok := options["playerId"].(string); ok {
		query["player.id"] = playerId
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListSeasonTeamPlayer(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterSeasonTeamPlayer(options, true)

	cur, err = r.Conn.Collection(r.seasonTeamPlayerCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListSeasonTeamPlayer Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountSeasonTeamPlayer(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterSeasonTeamPlayer(options, true)

	total, err := r.Conn.Collection(r.seasonTeamPlayerCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountSeasonTeamPlayer CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneSeasonTeamPlayer(ctx context.Context, options map[string]interface{}) (row *mongo_model.SeasonTeamPlayer, err error) {
	query, _ := generateQueryFilterSeasonTeamPlayer(options, false)

	err = r.Conn.Collection(r.seasonTeamPlayerCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneSeasonTeamPlayer FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneSeasonTeamPlayer(ctx context.Context, seasonTeamPlayer *mongo_model.SeasonTeamPlayer) (err error) {
	_, err = r.Conn.Collection(r.seasonTeamPlayerCollection).InsertOne(ctx, seasonTeamPlayer)
	if err != nil {
		logrus.Error("CreateSeasonTeamPlayer InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialSeasonTeamPlayer(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterSeasonTeamPlayer(options, false)

	_, err = r.Conn.Collection(r.seasonTeamPlayerCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialSeasonTeamPlayer UpdateOne:", err)
		return
	}

	return
}

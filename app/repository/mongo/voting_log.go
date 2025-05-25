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

func generateQueryFilterVotingLog(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if candidateId, ok := options["candidateId"].(string); ok {
		query["candidate.id"] = candidateId
	}
	if votingId, ok := options["votingId"].(string); ok {
		query["candidate.votingId"] = votingId
	}
	if memberId, ok := options["memberId"].(string); ok {
		query["memberId"] = memberId
	}

	return
}

func (r *mongoDbRepo) FetchListVotingLog(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterVotingLog(options, true)

	cur, err = r.Conn.Collection(r.votingLogCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListVotingLog Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountVotingLog(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterVotingLog(options, true)

	total, err := r.Conn.Collection(r.votingLogCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountVotingLog CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneVotingLog(ctx context.Context, options map[string]interface{}) (row *mongo_model.VotingLog, err error) {
	query, _ := generateQueryFilterVotingLog(options, false)

	err = r.Conn.Collection(r.votingLogCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneVotingLog FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneVotingLog(ctx context.Context, votingLog *mongo_model.VotingLog) (err error) {
	_, err = r.Conn.Collection(r.votingLogCollection).InsertOne(ctx, votingLog)
	if err != nil {
		logrus.Error("CreateVotingLog InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialVotingLog(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterVotingLog(options, false)

	_, err = r.Conn.Collection(r.votingLogCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialVotingLog UpdateOne:", err)
		return
	}

	return
}

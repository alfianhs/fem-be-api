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

func generateQueryFilterCandidate(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	if votingId, ok := options["votingId"].(string); ok {
		query["votingId"] = votingId
	}

	if seasonTeamPlayerId, ok := options["seasonTeamPlayerId"].(string); ok {
		query["seasonTeamPlayerId"] = seasonTeamPlayerId
	}

	return
}

func (r *mongoDbRepo) FetchListCandidate(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterCandidate(options, true)

	cur, err = r.Conn.Collection(r.candidateCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListCandidate Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountCandidate(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterCandidate(options, true)

	total, err := r.Conn.Collection(r.candidateCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountCandidate CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneCandidate(ctx context.Context, options map[string]interface{}) (row *mongo_model.Candidate, err error) {
	query, _ := generateQueryFilterCandidate(options, false)

	err = r.Conn.Collection(r.candidateCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneCandidate FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneCandidate(ctx context.Context, candidate *mongo_model.Candidate) (err error) {
	_, err = r.Conn.Collection(r.candidateCollection).InsertOne(ctx, candidate)
	if err != nil {
		logrus.Error("CreateCandidate InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialCandidate(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterCandidate(options, false)

	_, err = r.Conn.Collection(r.candidateCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialCandidate UpdateOne:", err)
		return
	}

	return
}

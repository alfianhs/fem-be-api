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

func generateQueryFilterVoting(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if status, ok := options["status"].(mongo_model.VotingStatus); ok {
		query["status"] = status
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchListVoting(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error) {
	query, findOptions := generateQueryFilterVoting(options, true)

	cur, err = r.Conn.Collection(r.votingCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchListVoting Find:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CountVoting(ctx context.Context, options map[string]interface{}) (total int64) {
	query, _ := generateQueryFilterVoting(options, true)

	total, err := r.Conn.Collection(r.votingCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountVoting CountDocuments:", err)
		return 0
	}

	return
}

func (r *mongoDbRepo) FetchOneVoting(ctx context.Context, options map[string]interface{}) (row *mongo_model.Voting, err error) {
	query, _ := generateQueryFilterVoting(options, false)

	err = r.Conn.Collection(r.votingCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneVoting FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneVoting(ctx context.Context, voting *mongo_model.Voting) (err error) {
	_, err = r.Conn.Collection(r.votingCollection).InsertOne(ctx, voting)
	if err != nil {
		logrus.Error("CreateVoting InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialVoting(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterVoting(options, false)

	_, err = r.Conn.Collection(r.votingCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialVoting UpdateOne:", err)
		return
	}

	return
}

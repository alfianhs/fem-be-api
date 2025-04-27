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

func generateQueryFilterMember(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if email, ok := options["email"].(string); ok {
		query["email"] = email
	}
	if emailToken, ok := options["emailToken"].(string); ok {
		query["emailToken"] = emailToken
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchOneMember(ctx context.Context, options map[string]interface{}) (row *mongo_model.Member, err error) {
	query, _ := generateQueryFilterMember(options, false)

	err = r.Conn.Collection(r.memberCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneMember FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneMember(ctx context.Context, member *mongo_model.Member) (err error) {
	_, err = r.Conn.Collection(r.memberCollection).InsertOne(ctx, member)
	if err != nil {
		logrus.Error("CreateMember InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialMember(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterMember(options, false)

	_, err = r.Conn.Collection(r.memberCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialMember UpdateOne:", err)
		return
	}

	return
}

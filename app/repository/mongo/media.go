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

func generateQueryFilterMedia(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchOneMedia(ctx context.Context, options map[string]interface{}) (row *mongo_model.Media, err error) {
	query, _ := generateQueryFilterMedia(options, false)

	err = r.Conn.Collection(r.mediaCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneMedia FindOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) CreateOneMedia(ctx context.Context, media *mongo_model.Media) (err error) {
	_, err = r.Conn.Collection(r.mediaCollection).InsertOne(ctx, media)
	if err != nil {
		logrus.Error("CreateMedia InsertOne:", err)
		return
	}
	return
}

func (r *mongoDbRepo) CreateManyMedia(ctx context.Context, medias []*mongo_model.Media) (err error) {
	docs := make([]interface{}, len(medias))
	for i, row := range medias {
		docs[i] = row
	}

	_, err = r.Conn.Collection(r.mediaCollection).InsertMany(ctx, docs)
	if err != nil {
		logrus.Error("CreateManyMedia InsertMany:", err)
		return
	}
	return
}

func (r *mongoDbRepo) UpdatePartialMedia(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterMedia(options, false)

	_, err = r.Conn.Collection(r.mediaCollection).UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdatePartialMedia UpdateOne:", err)
		return
	}

	return
}

func (r *mongoDbRepo) UpdateManyMediaPartial(ctx context.Context, options, field map[string]interface{}) (err error) {
	query, _ := generateQueryFilterMedia(options, false)

	_, err = r.Conn.Collection(r.mediaCollection).UpdateMany(ctx, query, bson.M{"$set": field})
	if err != nil {
		logrus.Error("UpdateManyMediaPartial UpdateMany:", err)
		return
	}

	return
}

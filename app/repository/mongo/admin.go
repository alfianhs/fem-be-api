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

func generateQueryFilterAdmin(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if username, ok := options["username"].(string); ok {
		query["username"] = username
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchOneAdmin(ctx context.Context, options map[string]interface{}) (row *mongo_model.Admin, err error) {
	query, _ := generateQueryFilterAdmin(options, false)

	err = r.Conn.Collection(r.adminCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneAdmin FindOne:", err)
		return
	}

	return
}

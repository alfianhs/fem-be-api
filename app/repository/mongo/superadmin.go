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

func generateQueryFilterSuperadmin(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// custom filter
	if email, ok := options["email"].(string); ok {
		query["email"] = email
	}

	return query, mongoOptions
}

func (r *mongoDbRepo) FetchOneSuperadmin(ctx context.Context, options map[string]interface{}) (row *mongo_model.Superadmin, err error) {
	query, _ := generateQueryFilterSuperadmin(options, false)

	err = r.Conn.Collection(r.superadminCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneSuperadmin FindOne:", err)
		return
	}

	return
}

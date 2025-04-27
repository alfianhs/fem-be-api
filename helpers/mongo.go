package helpers

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoDatabaseName(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	// get name
	dbName := strings.TrimPrefix(u.Path, "/")
	return dbName, nil
}

func ConnectMongoDB(timeout time.Duration, uri string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// client options mongo
	clientOptions := moptions.Client().ApplyURI(uri)

	// connect to mongo
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		panic(err)
	}

	// connection check
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging to MongoDB:", err)
		panic(err)
	}

	// get database name
	dbName, err := getMongoDatabaseName(uri)
	if err != nil {
		fmt.Println("Error getting database name:", err)
		panic(err)
	}

	return client.Database(dbName)
}

func CommonFilter(options map[string]any) map[string]any {
	query := map[string]any{
		"deletedAt": bson.M{
			"$eq": nil,
		},
	}

	if id, ok := options["id"].(primitive.ObjectID); ok {
		query["_id"] = id
	} else if id, ok := options["id"].(string); ok {
		obj, _ := primitive.ObjectIDFromHex(id)
		query["_id"] = obj
	}

	if ids, ok := options["ids"].([]primitive.ObjectID); ok {
		query["_id"] = bson.M{
			"$in": ids,
		}
	} else if ids, ok := options["ids"].([]string); ok {
		objIDs := make([]primitive.ObjectID, 0)
		for _, id := range ids {
			if obID, err := primitive.ObjectIDFromHex(strings.TrimSpace(id)); err == nil {
				objIDs = append(objIDs, obID)
			}
		}
		query["_id"] = bson.M{
			"$in": objIDs,
		}
	}

	return query
}

func CommonMongoFindOptions(options map[string]any) *moptions.FindOptions {
	// limit, offset & sort
	mongoOptions := moptions.Find()
	if offset, ok := options["offset"].(int64); ok {
		mongoOptions.SetSkip(offset)
	} else if offset, ok := options["offset"].(int); ok {
		mongoOptions.SetSkip(int64(offset))
	}

	if limit, ok := options["limit"].(int64); ok {
		mongoOptions.SetLimit(limit)
	} else if limit, ok := options["limit"].(int); ok {
		mongoOptions.SetLimit(int64(limit))
	}

	if sortBy, ok := options["sort"].(string); ok {
		sortDir, ok := options["dir"].(string)
		if !ok {
			sortDir = "asc"
		}

		sortQ := bson.D{}
		sortDirMongo := int(1)
		if strings.ToLower(sortDir) == "desc" {
			sortDirMongo = -1
		}
		sortQ = append(sortQ, bson.E{
			Key:   sortBy,
			Value: sortDirMongo,
		})
		mongoOptions.SetSort(sortQ)
	} else if sortBy, ok := options["sort"].(map[string]int); ok {
		sortQ := bson.D{}
		for k, sort := range sortBy {
			sortQ = append(sortQ, bson.E{
				Key:   k,
				Value: sort,
			})
		}
		mongoOptions.SetSort(sortQ)
	}

	if projection, ok := options["projection"].(map[string]int); ok {
		mongoOptions.SetProjection(projection)
	}

	return mongoOptions
}

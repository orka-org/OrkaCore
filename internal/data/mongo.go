package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoIndexes(ctx context.Context, db *mongo.Database) ([]string, error) {
	coll := db.Collection("users")
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "phone", Value: 1},
			},
			Options: options.Index().SetUnique(false),
		},
	}
	uniqueIndexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	res, err := coll.Indexes().CreateMany(ctx, append(indexModels, uniqueIndexModels...))
	if err != nil {
		return nil, err
	}

	return res, nil
}

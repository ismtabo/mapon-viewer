package helper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionHelper implements Mongo Collections interaction.
type CollectionHelper interface {
	// InsertOne inserts a document in Mongo Collection.
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	// FindOne finds a document in Mongo Collection by a given filter.
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
}

package storage

import (
	"context"
	"news-crawler/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoDBClient(uri, database, collection string) (*MongoDBClient, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoDBClient{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (mc *MongoDBClient) InsertMetadata(metadata *models.Metadata) error {
	collection := mc.client.Database(mc.database).Collection(mc.collection)
	_, err := collection.InsertOne(context.Background(), metadata)
	return err
}

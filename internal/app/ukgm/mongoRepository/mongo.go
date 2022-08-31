package mongorepository

import (
	"context"

	basicMongo "ukgm/internal/app/ukgm/mongoRepository/basicFunction"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

// MongoRepo implements facultative.Repository using MongoDB
type MongoRepo struct {
	client *mongo.Client
	db     *mongo.Database

	basicMongo.BasicMongo
}

// NewMongoRepo instantiates a new MongoRepo instance
func NewMongoRepo(url, dbName string) (*MongoRepo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		zap.S().Error()
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		zap.S().Error()
		return nil, err
	}

	result := &MongoRepo{
		client: client,
		db:     client.Database(dbName),
	}
	result.Init()
	return result, nil
}

// Collection returns collection by the given name
func (r *MongoRepo) Init() {
	r.BasicMongo = basicMongo.New(r.db)
	return
}

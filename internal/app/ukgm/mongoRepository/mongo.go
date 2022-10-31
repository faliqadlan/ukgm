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
	Client *mongo.Client
	Db     *mongo.Database

	basicMongo.BasicMongo
}

// NewMongoRepo instantiates a new MongoRepo instance
func NewMongoRepo(url, dbName string) (*MongoRepo, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		zap.S().Error()
		return nil, err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		zap.S().Error()
		return nil, err
	}

	result := &MongoRepo{
		Client: client,
		Db:     client.Database(dbName),
	}
	result.Init()
	return result, nil
}

// Collection returns collection by the given name
func (r *MongoRepo) Init() {
	r.BasicMongo = basicMongo.New(r.Db)
}

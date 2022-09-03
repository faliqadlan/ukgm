package basicMongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BasicMongo interface {
	Count(ctx context.Context, cname string, filter bson.M) (count int, err error)
	Distinct(ctx context.Context, cname, fieldName string, filter bson.M) (result []interface{}, err error)
	Aggregate(ctx context.Context, cname string, aggregatePipeline interface{}, options *options.AggregateOptions) (result *mongo.Cursor, err error)
	Find(ctx context.Context, cname string, filter bson.M, options *options.FindOptions) (result *mongo.Cursor, err error)
	Show(ctx context.Context, cname string, filter bson.M, option *options.FindOneOptions) (result *mongo.SingleResult)
	Store(ctx context.Context, cname string, data interface{}, options *options.InsertOneOptions) (result *mongo.InsertOneResult, err error)
	Update(ctx context.Context, cname string, filter bson.M, data interface{}, options *options.UpdateOptions) (result *mongo.UpdateResult, err error)
	UpdateMany(ctx context.Context, cname string, filter bson.M, data bson.M, options *options.UpdateOptions) (result *mongo.UpdateResult, err error)
	DeleteMany(ctx context.Context, cname string, filter bson.M, option *options.DeleteOptions) (result *mongo.DeleteResult, err error)
	Delete(ctx context.Context, cname string, filter bson.M, data interface{}, options *options.DeleteOptions) (result *mongo.DeleteResult, err error)
}

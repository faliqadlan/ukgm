package basicMongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	db *mongo.Database
}

func New(db *mongo.Database) (result *Repo) {
	result = &Repo{
		db: db,
	}
	return
}

func (r *Repo) Count(ctx context.Context, cname string, filter bson.M) (count int, err error) {
	collection := r.db.Collection(cname)

	count64, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}
	count = int(count64)

	return
}

func (r *Repo) Distinct(ctx context.Context, cname, fieldName string, filter bson.M) (result []interface{}, err error) {
	collection := r.db.Collection(cname)

	result, err = collection.Distinct(ctx, fieldName, filter)
	if err != nil {
		return
	}

	return
}

func (r *Repo) Aggregate(ctx context.Context, cname string, aggregatePipeline interface{}, options *options.AggregateOptions) (result *mongo.Cursor, err error) {
	collection := r.db.Collection(cname)

	result, err = collection.Aggregate(ctx, aggregatePipeline, options)

	return
}

func (r *Repo) Find(ctx context.Context, cname string, filter bson.M, options *options.FindOptions) (result *mongo.Cursor, err error) {
	collection := r.db.Collection(cname)

	result, err = collection.Find(ctx, filter, options)
	if err != nil {
		return
	}

	return
}

func (r *Repo) Show(ctx context.Context, cname string, filter bson.M, option *options.FindOneOptions) (result *mongo.SingleResult) {
	collection := r.db.Collection(cname)
	if option == nil {
		option = &options.FindOneOptions{}
	}
	result = collection.FindOne(ctx, filter, option)

	return
}

func (r *Repo) Store(ctx context.Context, cname string, data interface{}, options *options.InsertOneOptions) (result *mongo.InsertOneResult, err error) {
	collection := r.db.Collection(cname)
	result, err = collection.InsertOne(ctx, data, options)
	return
}

func (r *Repo) Update(ctx context.Context, cname string, filter bson.M, data interface{}, options *options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	collection := r.db.Collection(cname)
	result, err = collection.UpdateOne(ctx, filter, bson.M{
		"$set": data,
	}, options)

	return
}

func (r *Repo) UpdateMany(ctx context.Context, cname string, filter bson.M, data bson.M, options *options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	collection := r.db.Collection(cname)

	result, err = collection.UpdateMany(ctx, filter, bson.M{
		"$set": data,
	}, options)

	return
}

func (r *Repo) DeleteMany(ctx context.Context, cname string, filter bson.M, option *options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	collection := r.db.Collection(cname)

	result, err = collection.DeleteMany(ctx, filter, option)

	return
}

func (r *Repo) Delete(ctx context.Context, cname string, filter bson.M, data interface{}, options *options.DeleteOptions) (result *mongo.DeleteResult, err error) {
	var collection = r.db.Collection(cname)

	result, err = collection.DeleteOne(ctx, data, options)
	return
}

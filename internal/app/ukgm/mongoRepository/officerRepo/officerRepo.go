package officer_repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Repo struct {
	*mongorepository.MongoRepo
}

const Collection = "officers"

func (r *Repo) StoreData(data *entities.Officer, ctx context.Context) (result *entities.Officer, err error) {

	mongoResult, err := r.MongoRepo.Store(ctx, Collection, data, nil)

	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data.ID = utils.ObjectIDHex(mongoResult.InsertedID)
	return data, nil
}

func (r *Repo) UpdateData(data *entities.Officer, ctx context.Context) (err error) {
	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update officer")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = r.MongoRepo.Update(ctx, Collection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func (r *Repo) ShowData(filter bson.M, options *options.FindOneOptions, ctx context.Context) (result *entities.Officer, err error) {
	result = &entities.Officer{}
	err = r.Show(ctx, Collection, filter, options).Decode(result)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	s := string(data)
	message := fmt.Sprintf("Found a single document : " + s + "\n")
	r.MongoRepo.BasicMongo.LogPrint(message)

	return result, nil
}

func (r *Repo) FetchList(filter bson.M, opt *options.FindOptions, ctx context.Context) (result []entities.OfficerList, err error) {
	if opt == nil {
		opt = &options.FindOptions{}
	}

	cursor, err := r.MongoRepo.Find(ctx, Collection, filter, opt)
	if err != nil {
		zap.S().Error(err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		officer := entities.Officer{}
		if err = cursor.Decode(&officer); err != nil {
			zap.S().Error(err)
			return
		}
		result = append(result, officer.ToOfficerList())
	}

	return
}

package user_repo

import (
	"context"
	"encoding/json"
	"errors"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type UserRepo struct {
	*mongorepository.MongoRepo
}

const UserCollection = "users"

func (u *UserRepo) StoreData(ctx context.Context, data *entities.User) (result *entities.User, err error) {

	mongoResult, err := u.MongoRepo.Store(ctx, UserCollection, data, nil)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data.ID = utils.ObjectIDHex(mongoResult.InsertedID)
	result = data
	return
}

func (u *UserRepo) UpdateData(ctx context.Context, data *entities.User) (err error) {

	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update user")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = u.MongoRepo.Update(ctx, UserCollection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func (u *UserRepo) ShowUser(ctx context.Context, filter bson.M, options *options.FindOneOptions) (result *entities.User, err error) {
	result = &entities.User{}
	err = u.Show(context.TODO(), UserCollection, filter, options).Decode(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	_, err = json.Marshal(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return result, nil
}

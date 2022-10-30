package user_repo

import (
	"context"
	"encoding/json"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type UserRepo struct {
	*mongorepository.MongoRepo
}

const UserCollection = "users"

func (u *UserRepo) ShowUser(filter bson.M, options *options.FindOneOptions) (result *entities.User, err error) {
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

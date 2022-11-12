package user_service

import (
	"context"
	"errors"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	user_repo "ukgm/internal/app/ukgm/mongoRepository/userRepo"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type UserService struct {
	Repo *user_repo.UserRepo
}

const UserCollection = "users"

func New(repo *mongorepository.MongoRepo) *UserService {
	return &UserService{
		Repo: &user_repo.UserRepo{MongoRepo: repo},
	}
}

func (u *UserService) ShowUserByUid(uid string) (result *entities.User, err error) {
	var filter = bson.M{
		"timestamp.deleted_at": nil,
		"uid":                  uid,
	}

	result, err = u.Repo.ShowUser(context.TODO(), filter, nil)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("error in show user")
	}

	return result, nil
}

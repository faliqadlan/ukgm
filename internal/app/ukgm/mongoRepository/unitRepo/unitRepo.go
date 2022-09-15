package unit_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type UnitRepo struct {
	*mongorepository.MongoRepo
}

const UnitCollection = "units"

func (u *UnitRepo) StoreUnit(data *entities.Unit, ctx context.Context) (result *entities.Unit, err error) {
	mongoResult, err := u.Store(ctx, UnitCollection, data, nil)

	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data.ID = utils.ObjectIDHex(mongoResult.InsertedID)
	return data, nil
}

func (u *UnitRepo) ShowUnit(filter bson.M, options *options.FindOneOptions, ctx context.Context) (result *entities.Unit, err error) {
	result = &entities.Unit{}
	err = u.Show(ctx, UnitCollection, filter, options).Decode(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data, err := json.Marshal(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	s := string(data)
	message := fmt.Sprintf("Found a single document : " + s + "\n")
	u.MongoRepo.BasicMongo.LogPrint(message)

	return result, nil
}

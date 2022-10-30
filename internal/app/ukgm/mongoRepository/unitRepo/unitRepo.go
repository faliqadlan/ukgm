package unit_repo

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
	_, err = json.Marshal(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return result, nil
}

func (u *UnitRepo) FecthUnitList(filter bson.M, opt *options.FindOptions, ctx context.Context) (result []entities.UnitList, err error) {
	if opt == nil {
		opt = &options.FindOptions{}
	}

	cursor, err := u.Find(ctx, UnitCollection, filter, opt)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		unit := entities.Unit{}
		if err = cursor.Decode(&unit); err != nil {
			zap.S().Error(err)
			return nil, err
		}
		result = append(result, *unit.ToUnitList())
	}
	return result, nil
}

func (u *UnitRepo) UpdateUnit(data *entities.Unit, ctx context.Context) (err error) {
	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update unit")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = u.Update(ctx, UnitCollection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}

	return
}

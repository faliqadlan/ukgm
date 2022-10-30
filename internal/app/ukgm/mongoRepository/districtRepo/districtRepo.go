package districtrepo

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

type DistrictRepo struct {
	*mongorepository.MongoRepo
}

const DistrictCollection = "districts"

func (d *DistrictRepo) UpdateDistrict(data *entities.District, ctx context.Context) (err error) {
	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update district")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = d.Update(ctx, DistrictCollection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func (d *DistrictRepo) ShowDistrict(filter bson.M, options *options.FindOneOptions, ctx context.Context) (result *entities.District, err error) {
	result = &entities.District{}
	err = d.Show(ctx, DistrictCollection, filter, options).Decode(result)
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

func (d *DistrictRepo) FetchDistrictList(filter bson.M, opt *options.FindOptions, ctx context.Context) (result []entities.DistrictList, err error) {
	if opt == nil {
		opt = &options.FindOptions{}
	}

	cursor, err := d.Find(ctx, DistrictCollection, filter, opt)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		district := entities.District{}
		if err = cursor.Decode(&district); err != nil {
			zap.S().Error(err)
			return nil, err
		}
		result = append(result, *district.ToDistrictList())
	}
	return result, nil
}

func (d *DistrictRepo) FetchManyDistricts(filter bson.M, opt *options.FindOptions, ctx context.Context) (result []entities.District, err error) {
	if opt == nil {
		opt = &options.FindOptions{}
	}

	cursor, err := d.Find(ctx, DistrictCollection, filter, opt)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		district := entities.District{}
		if err = cursor.Decode(&district); err != nil {
			zap.S().Error(err)
			return nil, err
		}
		result = append(result, district)
	}
	return result, nil
}

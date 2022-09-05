package districtrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DistrictRepo struct {
	*mongorepository.MongoRepo
}

const DistrictCollection = "districts"

func (d *DistrictRepo) ShowDistrict(filter bson.M, options *options.FindOneOptions) (result *entities.District, err error) {
	result = &entities.District{}
	err = d.Show(context.TODO(), DistrictCollection, filter, options).Decode(result)
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
	d.MongoRepo.BasicMongo.LogPrint(message)

	return result, nil
}

func (d *DistrictRepo) FetchDistrictList(filter bson.M, opt *options.FindOptions) (result []entities.DistrictList, err error) {
	if opt == nil {
		opt = &options.FindOptions{}
	}

	cursor, err := d.Find(context.TODO(), DistrictCollection, filter, opt)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		district := entities.District{}
		if err = cursor.Decode(&district); err != nil {
			zap.S().Error(err)
			return nil, err
		}
		result = append(result, *district.ToDistrictList())
	}
	return result, nil
}

func (d *DistrictRepo) FetchManyDistricts(filter bson.M, opt *options.FindOptions) (result []entities.District, err error) {
	if opt == nil {
		opt = &options.FindOptions{}
	}

	cursor, err := d.Find(context.TODO(), DistrictCollection, filter, opt)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		district := entities.District{}
		if err = cursor.Decode(&district); err != nil {
			zap.S().Error(err)
			return nil, err
		}
		result = append(result, district)
	}
	return result, nil
}

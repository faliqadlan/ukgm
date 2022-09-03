package doctor_repo

import (
	"context"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DoctorRepo struct {
	*mongorepository.MongoRepo
}

const DoctorCollection = "doctors"

func (doctor *DoctorRepo) StoreDoctor(data *entities.Doctor) (result *entities.Doctor, err error) {

	mongoResult, err := doctor.Store(context.TODO(), DoctorCollection, data, nil)

	if err != nil {
		zap.S().Error(err)
		return &entities.Doctor{}, err
	}
	var insert = mongoResult.InsertedID
	data.ID = utils.ObjectIDHex(insert)
	return data, nil
}

func (doctor *DoctorRepo) ShowDoctor(filter bson.M, options *options.FindOneOptions) (result *entities.Doctor, err error) {
	result = &entities.Doctor{}
	err = doctor.Show(context.TODO(), DoctorCollection, filter, options).Decode(result)
	if err != nil {
		zap.S().Error(err)
		return &entities.Doctor{}, err
	}
	return result, nil
}

package diagnose_repo

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

type DiagnoseRepo struct {
	*mongorepository.MongoRepo
}

const DiagnoseCollection = "diagnoses"

func (d *DiagnoseRepo) StoreData(data *entities.Diagnose, ctx context.Context) (*entities.Diagnose, error) {
	mongoResult, err := d.Store(ctx, DiagnoseCollection, data, nil)

	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data.ID = utils.ObjectIDHex(mongoResult.InsertedID)
	return data, nil
}

func (d *DiagnoseRepo) ShowData(filter bson.M, options *options.FindOneOptions, ctx context.Context) (result *entities.Diagnose, err error) {
	result = &entities.Diagnose{}
	err = d.Show(ctx, DiagnoseCollection, filter, options).Decode(result)
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

func (d *DiagnoseRepo) UpdateData(data *entities.Diagnose, ctx context.Context) (err error) {
	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update diagnose")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = d.Update(ctx, DiagnoseCollection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

package unit_service

import (
	"context"
	"fmt"

	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	districtrepo "ukgm/internal/app/ukgm/mongoRepository/districtRepo"
	unit_repo "ukgm/internal/app/ukgm/mongoRepository/unitRepo"

	"firebase.google.com/go/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type UnitService struct {
	DistrictRepo *districtrepo.DistrictRepo
	UnitRepo     *unit_repo.UnitRepo
}

func New(repo *mongorepository.MongoRepo) *UnitService {
	return &UnitService{
		DistrictRepo: &districtrepo.DistrictRepo{MongoRepo: repo},
		UnitRepo:     &unit_repo.UnitRepo{MongoRepo: repo},
	}
}

func (u *UnitService) StoreData(data interface{}, auth *auth.Token) (result interface{}, err error) {

	errGeneral := fmt.Errorf("failed to store new unit")
	var ctx = context.TODO()
	var session mongo.Session

	var dataUnit = &entities.Unit{}

	dataUnit, ok := data.(*entities.Unit)
	if !ok {
		dataUnit = data.(*entities.Unit)
	}

	dataUnit.Timestamp.CreateNew(auth)

	districtId, err := primitive.ObjectIDFromHex(dataUnit.DistrictID)
	if err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  districtId,
	}

	district, err := u.DistrictRepo.ShowDistrict(filter, nil, context.TODO())
	if err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	filter = bson.M{
		"timestamp.deleted_at": nil,
		"name":                 dataUnit.Name,
		"district_id":          dataUnit.DistrictID,
	}

	_, errAvailable := u.UnitRepo.ShowUnit(filter, nil, context.TODO())
	if errAvailable == nil {
		zap.S().Error(err)
		return nil, fmt.Errorf("unit with name %v in district %v is already available", dataUnit.Name, district.Name)
	}

	if session, err = u.UnitRepo.Client.StartSession(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		_, err := u.UnitRepo.StoreUnit(dataUnit, sc)

		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		if err = session.CommitTransaction(sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.TODO()); abortErr != nil {
			zap.S().Error(err)
			return nil, errGeneral
		}
		zap.S().Error(err)
		return nil, errGeneral
	}

	return result, nil
}

func (u *UnitService) GetList(input ...string) (result interface{}, err error) {
	filter := bson.M{
		"timestamp.deleted_at": nil,
		"district_id":          input[0],
	}

	res, err := u.UnitRepo.FecthUnitList(filter, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}), context.TODO())
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return res, nil
}

// func (u *UnitService) GetDetail(input ...string) (result interface{}, err error) {

// 	var filter = bson.M{
// 		"timestamp.deleted_at": nil,
// 	}
// }

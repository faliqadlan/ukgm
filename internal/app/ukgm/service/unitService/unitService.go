package unit_service

import (
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	districtrepo "ukgm/internal/app/ukgm/mongoRepository/districtRepo"
	unit_repo "ukgm/internal/app/ukgm/mongoRepository/unitRepo"

	"firebase.google.com/go/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	dataUnit := data.(*entities.Unit)

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

	district, err := u.DistrictRepo.ShowDistrict(filter, nil)
	if err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	district.Units = append(district.Units, dataUnit)

	err = u.DistrictRepo.UpdateDistrict(district)
	if err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	return result, nil
}

package unit_service

import (
	"context"
	"fmt"
	"log"

	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	districtrepo "ukgm/internal/app/ukgm/mongoRepository/districtRepo"
	unit_repo "ukgm/internal/app/ukgm/mongoRepository/unitRepo"

	"firebase.google.com/go/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	errGeneral := fmt.Errorf("failed to store new unit")
	var ctx = context.Background()
	var client = u.UnitRepo.Client
	var session mongo.Session
	var collection = client.Database("ukgm_dev").Collection("units")

	dataUnit := data.(entities.Unit)

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

	district.Units = append(district.Units, &dataUnit)

	defer client.Disconnect(ctx)
	if session, err = client.StartSession(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		if err = u.DistrictRepo.UpdateDistrict(district); err != nil {
			log.Fatal(err)
		}

		resultUnit, err := collection.InsertOne(ctx, dataUnit)

		err1 := fmt.Errorf("error updating")
		if err1 != nil {
			log.Fatal(err, err1)
		}

		result = resultUnit

		if err = session.CommitTransaction(sc); err != nil {
			log.Fatal(err)
		}

		return nil
	}); err != nil {
		log.Println(err)
		// if err = session.AbortTransaction(ctx); err != nil {
		// 	zap.S().Error(err)
		// 	log.Fatal(err)
		// }
		// zap.S().Error(err)
		log.Fatal(err)
	}
	session.EndSession(ctx)

	return result, nil
}

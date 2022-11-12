package officer_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	officer_repo "ukgm/internal/app/ukgm/mongoRepository/officerRepo"
	user_repo "ukgm/internal/app/ukgm/mongoRepository/userRepo"

	firebase_api "ukgm/internal/app/config/firebase"

	"firebase.google.com/go/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type OfficerService struct {
	userRepo    *user_repo.UserRepo
	OfficerRepo *officer_repo.Repo
	firebaseApi *firebase_api.FirebaseApi
}

func New(repo *mongorepository.MongoRepo, firebaseApi *firebase_api.FirebaseApi) *OfficerService {
	return &OfficerService{
		OfficerRepo: &officer_repo.Repo{MongoRepo: repo},
		userRepo:    &user_repo.UserRepo{MongoRepo: repo},
		firebaseApi: firebaseApi,
	}
}

func (o *OfficerService) StoreData(data interface{}, auth *auth.Token) (result interface{}, err error) {
	errGeneral := fmt.Errorf("failed to store new office")
	var ctx = context.TODO()

	var dataOfficer = &entities.Officer{}
	var session mongo.Session
	dataOfficer, ok := data.(*entities.Officer)

	if !ok {
		return nil, errGeneral
	}

	if session, err = o.OfficerRepo.Client.StartSession(); err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		if dataOfficer.DistrictID == "" {
			dataOfficer.DistrictID = "kecamatan"
			dataOfficer.Level = "kecamatan"
		} else {
			dataOfficer.Level = "kelurahan"
		}

		dataOfficer.Timestamp.CreateNew(auth)

		res, err := o.OfficerRepo.StoreData(dataOfficer, sc)
		if err != nil {
			log.Println(err)
			zap.S().Error(err)
			return errGeneral
		}

		userData := entities.User{}

		userData.CreateNew(*auth, dataOfficer.DistrictID)

		result = res
		return
	}); err != nil {
		if abortErr := session.AbortTransaction(context.TODO()); abortErr != nil {
			zap.S().Error(err)
			return nil, errGeneral
		}
		zap.S().Error(err)
		return nil, errGeneral
	}

	
	return
}

func (o *OfficerService) SoftDelete(auth *auth.Token, input ...string) (err error) {
	var errGeneral = fmt.Errorf("failed to delete officer")
	officerId := input[0]

	oID, err := primitive.ObjectIDFromHex(officerId)
	if err != nil {
		zap.S().Error(err)
		return errors.New("failed to parse id")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oID,
	}

	show, err := o.OfficerRepo.ShowData(filter, nil, context.TODO())
	if err != nil {
		zap.S().Error(err)
		return errGeneral
	}

	show.Timestamp.Delete(auth)

	err = o.OfficerRepo.UpdateData(show, context.TODO())
	if err != nil {
		zap.S().Error(err)
		return errGeneral
	}

	return
}

func (o *OfficerService) GetList(input ...string) (result interface{}, err error) {
	var errGeneral = fmt.Errorf("failed to show officer list")
	filter := bson.M{
		"timestamp.deleted_at": nil,
		"level":                "kelurahan",
	}

	if input[0] != "" {
		filter["district_id"] = input[0]
	}

	switch {
	case input[0] != "":
		filter["district_id"] = input[0]
	case input[1] != "":
		filter["profesion"] = input[1]
	}

	res, err := o.OfficerRepo.FetchList(filter, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}), context.TODO())
	if err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if input[2] == "true" {
		filter = bson.M{
			"timestamp.deleted_at": nil,
			"district_id":          "kecamatan",
			"level":                "kecamatan",
		}

		if input[1] != "" {
			filter["profesion"] = input[1]
		}

		res1, err := o.OfficerRepo.FetchList(filter, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}), context.TODO())
		if err != nil {
			zap.S().Error(err)
			return nil, errGeneral
		}

		res = append(res, res1...)
	}

	result = res
	return
}

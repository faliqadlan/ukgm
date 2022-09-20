package patient_service

import (
	"context"
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	districtrepo "ukgm/internal/app/ukgm/mongoRepository/districtRepo"
	patient_repo "ukgm/internal/app/ukgm/mongoRepository/patientRepo"
	unit_repo "ukgm/internal/app/ukgm/mongoRepository/unitRepo"
	"ukgm/internal/app/utils"

	"firebase.google.com/go/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type PatientService struct {
	DistrictRepo *districtrepo.DistrictRepo
	UnitRepo     *unit_repo.UnitRepo
	PatientRepo  *patient_repo.PatientRepo
}

func New(repo *mongorepository.MongoRepo) *PatientService {
	return &PatientService{
		DistrictRepo: &districtrepo.DistrictRepo{MongoRepo: repo},
		UnitRepo:     &unit_repo.UnitRepo{MongoRepo: repo},
		PatientRepo:  &patient_repo.PatientRepo{MongoRepo: repo},
	}
}

func (p *PatientService) StoreData(data interface{}, auth *auth.Token) (result interface{}, err error) {
	errGeneral := fmt.Errorf("failed to store new patient")
	var ctx = context.Background()
	var session mongo.Session
	var dataPat = &entities.Patient{}

	dataPat, ok := data.(*entities.Patient)
	if !ok {
		return nil, errGeneral
	}

	dataPat.Timestamp.CreateNew(auth)

	if session, err = p.PatientRepo.Client.StartSession(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		districtId, err := primitive.ObjectIDFromHex(dataPat.DistrictID)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		filter := bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  districtId,
		}

		district, err := p.DistrictRepo.ShowDistrict(filter, nil, context.TODO())
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		unitId, err := primitive.ObjectIDFromHex(dataPat.UnitID)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		filter = bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  unitId,
		}

		unit, err := p.UnitRepo.ShowUnit(filter, nil, context.TODO())
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		dataPat.NameUnit = unit.Name
		resultPat, err := p.PatientRepo.StorePatient(dataPat, sc)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		unit.Patients = append(unit.Patients, resultPat)

		if err = p.UnitRepo.UpdateUnit(unit, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		for i, unitD := range district.Units {
			if utils.ObjectIDHex(unitD.ID) == utils.ObjectIDHex(dataPat.UnitID) {
				district.Units[i].Patients = append(district.Units[i].Patients, resultPat)
				break
			}
		}

		if err = p.DistrictRepo.UpdateDistrict(district, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}
		// err = p.DistrictRepo.UpdateDistrict(district, sc)
		// err = fmt.Errorf("unable to update district: %v", err)
		// if err != nil {
		// 	zap.S().Error(err)
		// 	return errGeneral
		// }

		if err = session.CommitTransaction(sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		result = resultPat
		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			zap.S().Error(err)
			return nil, errGeneral
		}
		zap.S().Error(err)
		return nil, errGeneral
	}

	return result, nil
}

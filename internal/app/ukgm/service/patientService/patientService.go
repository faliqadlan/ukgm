package patient_service

import (
	"context"
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
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

	mrNumber, err := p.PatientRepo.GenerateMRNumber(ctx)
	if err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}
	dataPat.MrNumber = mrNumber

	dataPat.Timestamp.CreateNew(auth)

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"name":                 dataPat.Name,
		"dob":                  dataPat.Dob,
	}

	_, errAvailable := p.PatientRepo.Showpatient(filter, nil, ctx)
	if errAvailable == nil {
		zap.S().Error(err)
		return nil, fmt.Errorf("patient with name %v and date of birth %v is already available", dataPat.Name, dataPat.Dob.Format("2006-01-02"))
	}

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

		mapPatient := make(map[string]*entities.Patient)
		mapPatient[utils.ObjectIDHex(resultPat.ID)] = resultPat

		if unit.Patients == nil {
			unit.Patients = mapPatient
		} else {
			unit.Patients[utils.ObjectIDHex(resultPat.ID)] = resultPat
		}

		if err = p.UnitRepo.UpdateUnit(unit, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		if district.Units[resultPat.UnitID].Patients == nil {
			district.Units[resultPat.UnitID].Patients = mapPatient
		} else {
			district.Units[resultPat.UnitID].Patients[utils.ObjectIDHex(resultPat.ID)] = resultPat
		}

		if err = p.DistrictRepo.UpdateDistrict(district, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		// for test session error
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

func (p *PatientService) FetchIndex(input ...interface{}) (result interface{}, count int, err error) {

	pr := input[0].(*request.PageRequest)

	res, count, err := p.PatientRepo.FetchIndexPatient(context.TODO(), pr)
	if err != nil {
		return nil, count, err
	}

	return res, count, nil
}

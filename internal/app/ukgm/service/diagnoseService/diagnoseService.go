package diagnose_service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	diagnose_repo "ukgm/internal/app/ukgm/mongoRepository/diagnoseRepo"
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

type DiagnoseService struct {
	DistrictRepo *districtrepo.DistrictRepo
	UnitRepo     *unit_repo.UnitRepo
	PatientRepo  *patient_repo.PatientRepo
	DiagnoseRepo *diagnose_repo.DiagnoseRepo
}

func New(repo *mongorepository.MongoRepo) *DiagnoseService {
	return &DiagnoseService{
		DistrictRepo: &districtrepo.DistrictRepo{MongoRepo: repo},
		UnitRepo:     &unit_repo.UnitRepo{MongoRepo: repo},
		PatientRepo:  &patient_repo.PatientRepo{MongoRepo: repo},
		DiagnoseRepo: &diagnose_repo.DiagnoseRepo{MongoRepo: repo},
	}
}

func (d *DiagnoseService) StoreData(data interface{}, auth *auth.Token) (result interface{}, err error) {
	var errGeneral = fmt.Errorf("failed to store new Diagnose")
	var ctx = context.Background()
	var session mongo.Session
	var dataStore = &entities.Diagnose{}

	dataInput := data.([]interface{})
	dataStore, ok := dataInput[0].(*entities.Diagnose)
	if !ok {
		return nil, errGeneral
	}

	dataStore.Timestamp.CreateNew(auth)

	replace, ok := dataInput[1].(bool)
	if !ok {
		return nil, errGeneral
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, errGeneral
	}

	filter := bson.M{
		"patient_id":           dataStore.PatientID,
		"timestamp.deleted_at": nil,
		"timestamp.updated_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(utils.StartOfMonth(time.Now().In(location).Local())),
			"$lte": primitive.NewDateTimeFromTime(utils.EndOfMonth(time.Now().In(location).Local())),
		},
	}
	resAvailable, errAvailable := d.DiagnoseRepo.ShowData(filter, nil, ctx)
	zap.S().Error(errAvailable)
	if errAvailable == nil {
		if !replace {
			return nil, fmt.Errorf("this patient have diagnose data for %v created at %v, are you sure to replace it ?", resAvailable.Timestamp.CreatedAt.Month().String(), resAvailable.Timestamp.CreatedAt.Format("2006-01-02"))
		} else {
			resAvailable.Timestamp.Delete(auth)
			err := d.DiagnoseRepo.UpdateData(resAvailable, ctx)
			if err != nil {
				return nil, errGeneral
			}
		}
	}

	if session, err = d.DiagnoseRepo.Client.StartSession(); err != nil {
		zap.S().Error(err)
		return nil, err
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		districtId, err := primitive.ObjectIDFromHex(dataStore.DistrictID)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		filter := bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  districtId,
		}

		district, err := d.DistrictRepo.ShowDistrict(filter, nil, context.TODO())
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		unitId, err := primitive.ObjectIDFromHex(dataStore.UnitID)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		filter = bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  unitId,
		}

		unit, err := d.UnitRepo.ShowUnit(filter, nil, context.TODO())
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		patientId, err := primitive.ObjectIDFromHex(dataStore.PatientID)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		filter = bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  patientId,
		}

		patient, err := d.PatientRepo.Showpatient(filter, nil, context.TODO())
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		// set doctor
		if district.Doctor != nil {
			for _, v := range district.Doctor {
				dataStore.DoctorName = v.Name
			}
		}

		// set nurse
		if district.Nurse != nil {
			for _, v := range district.Nurse {
				dataStore.NurseName = v.Name
			}
		}

		resultStore, err := d.DiagnoseRepo.StoreData(dataStore, sc)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		patient.LastDiagnose = dataStore

		if err = d.PatientRepo.UpdatePatient(patient, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		unit.Patients[resultStore.PatientID].LastDiagnose = resultStore
		if err = d.UnitRepo.UpdateUnit(unit, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		district.Units[resultStore.UnitID].Patients[resultStore.PatientID].LastDiagnose = resultStore
		if err = d.DistrictRepo.UpdateDistrict(district, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		// for test session error
		// err = fmt.Errorf("unable to update district: %v", err)
		// if err != nil {
		// 	zap.S().Error(err)
		// 	return errGeneral
		// }

		if err = session.CommitTransaction(sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		result = resultStore
		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			zap.S().Error(err)
			return nil, errGeneral
		}
		zap.S().Error(err)
		return nil, errGeneral
	}
	return
}

func (d *DiagnoseService) GetDetail(input ...string) (result interface{}, err error) {
	var errGeneral = fmt.Errorf("failed to show diagnose")
	patientId := input[0]

	oID, err := primitive.ObjectIDFromHex(patientId)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("failed to parse id")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oID,
	}

	res, err := d.PatientRepo.Showpatient(filter, nil, context.TODO())
	if err != nil {
		return nil, errGeneral
	}
	result = res.LastDiagnose
	return
}

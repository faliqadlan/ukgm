package patient_service

import (
	"context"
	"errors"
	"fmt"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	diagnose_repo "ukgm/internal/app/ukgm/mongoRepository/diagnoseRepo"
	districtrepo "ukgm/internal/app/ukgm/mongoRepository/districtRepo"
	patient_repo "ukgm/internal/app/ukgm/mongoRepository/patientRepo"
	unit_repo "ukgm/internal/app/ukgm/mongoRepository/unitRepo"

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
	DiagnoseRepo *diagnose_repo.DiagnoseRepo
}

func New(repo *mongorepository.MongoRepo) *PatientService {
	return &PatientService{
		DistrictRepo: &districtrepo.DistrictRepo{MongoRepo: repo},
		UnitRepo:     &unit_repo.UnitRepo{MongoRepo: repo},
		PatientRepo:  &patient_repo.PatientRepo{MongoRepo: repo},
		DiagnoseRepo: &diagnose_repo.DiagnoseRepo{MongoRepo: repo},
	}
}

func (p *PatientService) StoreData(data interface{}, auth *auth.Token) (result interface{}, err error) {
	errGeneral := fmt.Errorf("failed to store new patient")
	var ctx = context.TODO()
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
		if abortErr := session.AbortTransaction(context.TODO()); abortErr != nil {
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

func (p *PatientService) UpdateData(data interface{}, auth *auth.Token, input ...string) (result interface{}, err error) {

	errGeneral := fmt.Errorf("failed to update patient")
	var ctx = context.TODO()
	var session mongo.Session
	patientId := input[0]

	oID, err := primitive.ObjectIDFromHex(patientId)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("failed to parse id")
	}

	dataReq, ok := data.(*entities.PatientUpdateRequest)
	if !ok {
		return nil, errGeneral
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

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		filter := bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  oID,
		}

		show, err := p.PatientRepo.Showpatient(filter, nil, sc)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		patReq, err := dataReq.ToEntity()
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		show.Name = patReq.Name
		show.Dob = patReq.Dob
		show.NoHp = patReq.NoHp
		show.Gender = patReq.Gender
		show.Timestamp.Update(auth)

		if err := p.PatientRepo.UpdatePatient(show, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		result = show

		if err = session.CommitTransaction(sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}
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

func (p *PatientService) SoftDelete(auth *auth.Token, input ...string) (err error) {
	var errGeneral = fmt.Errorf("failed to delete patient")
	var ctx = context.TODO()
	var session mongo.Session
	patientId := input[0]

	oID, err := primitive.ObjectIDFromHex(patientId)
	if err != nil {
		zap.S().Error(err)
		return errors.New("failed to parse id")
	}

	if session, err = p.PatientRepo.Client.StartSession(); err != nil {
		zap.S().Error(err)
		return errGeneral
	}

	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) (err error) {

		filter := bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  oID,
		}

		show, err := p.PatientRepo.Showpatient(filter, nil, sc)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		show.Timestamp.Delete(auth)
		if show.LastDiagnose != nil {
			show.LastDiagnose.Timestamp.Delete(auth)
		}

		if err := p.PatientRepo.UpdatePatient(show, sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		filter = bson.M{
			"timestamp.deleted_at": nil,
			"patient_id":           patientId,
		}

		diagnoses, err := p.DiagnoseRepo.FetchDiagnoseList(sc, filter, nil)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		if len(diagnoses) > 0 {
			for i := range diagnoses {
				diagnoses[i].Timestamp.Delete(auth)

				if err := p.DiagnoseRepo.UpdateData(&diagnoses[i], sc); err != nil {
					zap.S().Error(err)
					return errGeneral
				}
			}
		}

		if err = session.CommitTransaction(sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}
		return
	}); err != nil {
		if abortErr := session.AbortTransaction(context.TODO()); abortErr != nil {
			zap.S().Error(err)
			return errGeneral
		}
		zap.S().Error(err)
		return errGeneral
	}

	return
}

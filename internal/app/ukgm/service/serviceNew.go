package service

import (
	"fmt"
	"net/url"

	"ukgm/internal/app/config/env"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	diagnose_service "ukgm/internal/app/ukgm/service/diagnoseService"
	district_service "ukgm/internal/app/ukgm/service/district_service"
	doctor_service "ukgm/internal/app/ukgm/service/doctorService"
	officer_service "ukgm/internal/app/ukgm/service/officerService"
	patient_service "ukgm/internal/app/ukgm/service/patientService"
	unit_service "ukgm/internal/app/ukgm/service/unitService"
	user_service "ukgm/internal/app/ukgm/service/userService"

	firebase_api "ukgm/internal/app/config/firebase"

	"go.uber.org/zap"
)

type ServiceNew struct {
	Repo *mongorepository.MongoRepo
	*doctor_service.DoctorService
	*district_service.DistrictService
	user_service.UserServiceColl
	UnitSrvStoreData     StoreDataSrv
	UnitSrvGetList       GetListsrv
	PatientSrvStoreData  StoreDataSrv
	PatientSrvFetchIndex FetchIndexSrv
	PatientSrvUpdateData UpdateDataSrv
	PatientSrvSoftDelete DeleteDataSrv
	DiagnoseSrvStoreData StoreDataSrv
	DiagnoseSrvGetDetail GetDetailSrv
	DiagnoseSrvGetList   GetListsrv
	OfficerSrvStoreData  StoreDataSrv
	OfficerSrvSoftDelete DeleteDataSrv
	OfficerSrvGetList    GetListsrv
}

func ServiceNewDeclare(firebaseApi *firebase_api.FirebaseApi) (result *ServiceNew, err error) {
	repo, err := MongoDeclare()
	if err != nil {
		zap.S().Error(err)
		return
	}

	result = &ServiceNew{
		Repo: repo,
	}
	result.Init(firebaseApi)
	return
}

func MongoDeclare() (repo *mongorepository.MongoRepo, err error) {
	dbUsername := env.Get("DB_USERNAME")
	dbPassword := env.Get("DB_PASSWORD")
	dbHost := env.Get("DB_HOST")
	dbPort := env.Get("DB_PORT")
	dbName := env.Get("DB_DATABASE")
	dbURL := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", dbUsername, url.QueryEscape(dbPassword), dbHost, dbPort)
	repo, err = mongorepository.NewMongoRepo(dbURL, dbName)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func (s *ServiceNew) Init(firebaseApi *firebase_api.FirebaseApi) {
	s.DoctorService = doctor_service.New(s.Repo)
	s.DistrictService = district_service.New(s.Repo)
	s.UserServiceColl = user_service.New(s.Repo)
	s.UnitSrvStoreData = unit_service.New(s.Repo)
	s.UnitSrvGetList = unit_service.New(s.Repo)
	s.PatientSrvStoreData = patient_service.New(s.Repo)
	s.PatientSrvFetchIndex = patient_service.New(s.Repo)
	s.PatientSrvUpdateData = patient_service.New(s.Repo)
	s.PatientSrvSoftDelete = patient_service.New(s.Repo)
	s.DiagnoseSrvStoreData = diagnose_service.New(s.Repo)
	s.DiagnoseSrvGetDetail = diagnose_service.New(s.Repo)
	s.DiagnoseSrvGetList = diagnose_service.New(s.Repo)
	s.OfficerSrvStoreData = officer_service.New(s.Repo, firebaseApi)
	s.OfficerSrvSoftDelete = officer_service.New(s.Repo, firebaseApi)
	s.OfficerSrvGetList = officer_service.New(s.Repo, firebaseApi)
}

func MongoDeclareForTesting() (repo *mongorepository.MongoRepo, err error) {
	dbUsername := env.Get("DB_USERNAME")
	dbPassword := env.Get("DB_PASSWORD")
	dbHost := env.Get("DB_HOST")
	dbPort := env.Get("DB_PORT")
	dbName := env.Get("DB_DATABASE")
	dbURL := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", dbUsername, url.QueryEscape(dbPassword), dbHost, dbPort)
	repo, err = mongorepository.NewMongoRepo(dbURL, dbName)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

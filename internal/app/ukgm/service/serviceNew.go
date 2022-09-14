package service

import (
	"fmt"
	"net/url"

	"ukgm/internal/app/config/env"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	district_service "ukgm/internal/app/ukgm/service/district_service"
	doctor_service "ukgm/internal/app/ukgm/service/doctorService"
	user_service "ukgm/internal/app/ukgm/service/userService"

	"go.uber.org/zap"
)

type ServiceNew struct {
	Repo *mongorepository.MongoRepo
	*doctor_service.DoctorService
	*district_service.DistrictService
	user_service.UserServiceColl
}

func ServiceNewDeclare() (result *ServiceNew, err error) {
	repo, err := MongoDeclare()
	if err != nil {
		zap.S().Error(err)
		return
	}

	result = &ServiceNew{
		Repo: repo,
	}
	result.Init()
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

func (s *ServiceNew) Init() {
	s.DoctorService = doctor_service.New(s.Repo)
	s.DistrictService = district_service.New(s.Repo)
	s.UserServiceColl = user_service.New(s.Repo)
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

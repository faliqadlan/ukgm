package service

import (
	"fmt"
	"net/url"

	"ukgm/internal/app/config/env"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"

	"go.uber.org/zap"
)

type ServiceNew struct {
	Repo *mongorepository.MongoRepo
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

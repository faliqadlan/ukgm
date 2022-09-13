package user_service

import (
	"fmt"
	"net/url"
	"testing"
	"ukgm/internal/app/config/env"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

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

func TestShowUserByUid(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(repo)

	srv := New(repo)

	t.Run("success show user by uid ", func(t *testing.T) {
		res, err := srv.ShowUserByUid("n8xgf00q7udFuCX7jbTIhtyf2gk2")
		if err != nil {
			t.Error(err)
		}
		t.Log(res)
		assert.Equal(t, err, nil)
	})
}

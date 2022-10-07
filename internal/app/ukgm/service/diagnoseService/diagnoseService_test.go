package diagnose_service

import (
	"fmt"
	"net/url"
	"testing"
	"time"
	"ukgm/internal/app/config"
	"ukgm/internal/app/config/env"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"

	"firebase.google.com/go/auth"
	"github.com/go-playground/assert"
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

func init() {
	config.SetupLogger()
}

func TestStoreData(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)
	t.Log(srv)

	t.Run("success store data", func(t *testing.T) {
		data := &entities.Diagnose{
			PatientID:  "63328dfd6f67b4f7e7fe2b82",
			UnitID:     "632914b5b3927d861c734ac5",
			DistrictID: "631dd39ef9c8ba488679be28",
			Note:       time.Now().String(),
			Condition:  ukgm.KondisiBaik,
		}

		dataInput := []interface{}{data, true}

		res, err := srv.StoreData(dataInput, &auth.Token{})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(res)
	})

	t.Run("not replace", func(t *testing.T) {
		data := &entities.Diagnose{
			PatientID:  "63328dfd6f67b4f7e7fe2b82",
			UnitID:     "632914b5b3927d861c734ac5",
			DistrictID: "631dd39ef9c8ba488679be28",
			Note:       time.Now().String(),
			Condition:  ukgm.KondisiBaik,
		}

		dataInput := []interface{}{data, false}

		res, err := srv.StoreData(dataInput, &auth.Token{})
		t.Log(err)
		assert.NotEqual(t, err, nil)
		t.Log(res)
	})

	t.Run("fail session", func(t *testing.T) {
		data := &entities.Diagnose{
			PatientID:  "63328dfd6f67b4f7e7fe2b82",
			UnitID:     "632914b5b3927d861c734ac5",
			DistrictID: "631dd39ef9c8ba488679be28",
			Note:       time.Now().String(),
			Condition:  ukgm.KondisiBaik,
		}

		dataInput := []interface{}{data, false}

		res, err := srv.StoreData(dataInput, &auth.Token{})
		t.Log(err)
		assert.NotEqual(t, err, nil)
		t.Log(res)
	})
}

func TestGetDetail(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)
	t.Log(srv)

	t.Run("success get data", func(t *testing.T) {
		res, err := srv.GetDetail("63328dfd6f67b4f7e7fe2b82")
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(res)
	})
}

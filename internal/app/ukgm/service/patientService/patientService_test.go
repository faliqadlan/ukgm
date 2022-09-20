package patient_service

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"ukgm/internal/app/config/env"
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

func TestStorePatirnt(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	t.Run("failed store patient", func(t *testing.T) {
		patient := &entities.Patient{
			DistrictID: "631dd39ef9c8ba488679be28",
			UnitID:     "6323169e850369d097c4fd52",
			Name:       "patient 1",
		}

		res, err := srv.StoreData(patient, &auth.Token{})
		t.Log(err)
		log.Println(err)
		assert.NotEqual(t, err, nil)
		t.Log(res)
	})

	t.Run("success store patient", func(t *testing.T) {
		patient := &entities.Patient{
			DistrictID: "631dd39ef9c8ba488679be28",
			UnitID:     "6323169e850369d097c4fd53",
			Name:       "patient tiga",
		}

		res, err := srv.StoreData(patient, &auth.Token{})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(res)
	})
}

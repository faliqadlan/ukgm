package patient_service

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"
	"ukgm/internal/app/config"
	"ukgm/internal/app/config/env"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

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

func TestStorePatient(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)
	t.Log(srv)

	loc, err := time.LoadLocation("Asia/Jakarta")
	t.Log(err)

	dob := time.Now().Local().In(loc)
	dob = utils.StartOfDay(dob)
	t.Log(dob)
	t.Log(dob.Location())

	t.Run("failed store patient", func(t *testing.T) {
		patient := &entities.Patient{
			DistrictID: "631dd39ef9c8ba488679be28",
			UnitID:     "632914b5b3927d861c734ac4",
			Name:       "patient 1",
			Dob:        &dob,
			Gender:     "netral",
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
			UnitID:     "632914b5b3927d861c734ac5",
			Name:       "budi",
			Dob:        &dob,
			Gender:     "netral",
		}

		res, err := srv.StoreData(patient, &auth.Token{})
		t.Log(err)
		assert.Equal(t, err, nil)
		t.Log(res)
	})
}

func TestFetchIndexPatient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	t.Run("success", func(t *testing.T) {

		pgReq := &request.PageRequest{
			// Search: "satu",
			// PerPage:  10,
			// Paginate: 1,
			// Page:     1,
			Sorts: []request.Sort{{Value: "asc", Option: "name"}},
		}

		res, count, err := srv.FetchIndex(pgReq)
		assert.Equal(t, err, nil)
		t.Log(count)
		t.Log(res)

	})
}

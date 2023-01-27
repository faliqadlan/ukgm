package patient_service

import (
	"context"
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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			UnitID:     "632914a7ff406c851d8ec418",
			Name:       "budi d",
			Dob:        &dob,
			Gender:     "P",
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

		pgReq.Filters = append(pgReq.Filters, request.Filter{Option: "unit_id", Operator: "=", Value: "632914a7ff406c851d8ec418"})

		res, count, err := srv.FetchIndex(pgReq)
		assert.Equal(t, err, nil)
		t.Log(count)
		t.Log(res)

	})
}

func TestUpdatePatient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	t.Run("success", func(t *testing.T) {

		oid, err := primitive.ObjectIDFromHex("634eb49f4e6e17cd1dcc5836")
		if err != nil {
			t.Error(err)
		}

		filter := bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  oid,
		}

		data, err := srv.PatientRepo.Showpatient(filter, nil, context.TODO())
		if err != nil {
			t.Error(err)
		}

		req := &entities.PatientUpdateRequest{
			Name:   data.Name,
			Gender: data.Gender,
			NoHp:   data.NoHp,
			Dob:    "1999-05-05",
		}

		res, err := srv.UpdateData(req, &auth.Token{}, "634eb49f4e6e17cd1dcc5836")
		if err != nil {
			t.Error(err)
		}

		t.Log(res)
	})
}

func TestDeledePatient(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	t.Run("success", func(t *testing.T) {
		patientId := "63be1768ede33d6221da5afe"

		err := srv.SoftDelete(&auth.Token{}, patientId)

		if err != nil {
			t.Log(err)
		}
	})
}

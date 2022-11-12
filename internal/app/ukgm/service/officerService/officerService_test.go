package officer_service

import (
	"fmt"
	"log"
	"net/url"
	"testing"
	"ukgm/internal/app/config"
	"ukgm/internal/app/config/env"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"firebase.google.com/go/auth"
	"github.com/go-playground/assert"
	"github.com/google/uuid"
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

/*
	Cipete Selatan
	Gandaria Selatan
	Cilandak Barat
	Lebak Bulus
	Pondok Labu
	test kelurahan
*/

var districtNameList = []string{"Cipete Selatan", "Gandaria Selatan", "Cilandak Barat", "Lebak Bulus", "Pondok Labu", "Test kelurahan"}

// dev
// var districtIdList = []string{"631dd39ef9c8ba488679be23", "631dd39ef9c8ba488679be24", "631dd39ef9c8ba488679be25", "631dd39ef9c8ba488679be26", "631dd39ef9c8ba488679be27", "631dd39ef9c8ba488679be28"}

// live
var districtIdList = []string{"6319e9dbb0aa169ea947cf63", "6319e9dbb0aa169ea947cf64", "6319e9dbb0aa169ea947cf65", "6319e9dbb0aa169ea947cf66", "6319e9dcb0aa169ea947cf67", "6319e9dcb0aa169ea947cf68"}

var doctoNameDistrictList = []string{"drg. Ira Mustika", "drg. Lilie Prima Likuningsih", "drg. Hany Faridah", "drg. Eva Susiyanti", "drg.Dian Kemala Sari", "dokter test kelurahan"}

var nurseNameDistrictList = []string{"Rilia Ika Puspita", "Fertiyan Dwi Cahyo", "Lucky Widyawati", "Masyita", "ina Fajarwati", "perawat test kelurahan"}

var doctorNameKecamatanList = []string{"drg. Shirly Sungkar", "drg. Rebecca Grace Agustina"}

var nurseNameKecamatanList = []string{"Herlina Novi", "Rahma chairunnisa"}

func init() {
	config.SetupLogger()
}

func TestInputOfficer(t *testing.T) {

	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo, )

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println(len(districtIdList), len(doctoNameDistrictList), len(nurseNameDistrictList), len(nurseNameKecamatanList), len(doctorNameKecamatanList))

	// input dokter kelurahan
	for i := 0; i < len(districtIdList); i++ {

		officerData := &entities.Officer{
			DistrictID:   districtIdList[i],
			Name:         doctoNameDistrictList[i],
			Level:        "kelurahan",
			Profesion:    "doctor",
			DistrictName: districtNameList[i],
		}

		res, err := srv.StoreData(officerData, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)
	}

	// input nurse kelurahan
	for i := 0; i < len(districtIdList); i++ {
		officerData := &entities.Officer{
			DistrictID:   districtIdList[i],
			Name:         nurseNameDistrictList[i],
			Level:        "kelurahan",
			Profesion:    "nurse",
			DistrictName: districtNameList[i],
		}

		res, err := srv.StoreData(officerData, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)
	}

	// input doctor kecamatan
	for i := 0; i < len(doctorNameKecamatanList); i++ {
		officerData := &entities.Officer{
			Name:      doctorNameKecamatanList[i],
			Profesion: "doctor",
		}

		res, err := srv.StoreData(officerData, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)
	}

	// input nurse kecamatan
	for i := 0; i < len(nurseNameKecamatanList); i++ {
		officerData := &entities.Officer{
			Name:      nurseNameKecamatanList[i],
			Profesion: "nurse",
		}

		res, err := srv.StoreData(officerData, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)
	}
}

func TestNurse(t *testing.T) {
	// 6357f6aa495cc77b8bd6e50c

	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	t.Run("store nurse and delete", func(t *testing.T) {
		data := &entities.Officer{
			DistrictID: "6319e9dcb0aa169ea947cf68",
			Name:       "perawat test kelurahan " + uuid.NewString(),
			Profesion:  "nurse",
		}

		res, err := srv.StoreData(data, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)

		delData := res.(*entities.Officer)

		err = srv.SoftDelete(&auth.Token{}, utils.ObjectIDHex(delData.ID))
		assert.Equal(t, err, nil)
	})

	t.Run("fail soft test delete", func(t *testing.T) {
		err := srv.SoftDelete(&auth.Token{}, "6357f6aa495cc77b8bd6e50c")
		assert.NotEqual(t, err, nil)
	})

	t.Run("show list", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "nurse", "true")
		assert.Equal(t, err, nil)
		log.Println(res)
	})

	t.Run("store nurse", func(t *testing.T) {
		data := &entities.Officer{
			DistrictID: "6319e9dcb0aa169ea947cf68",
			Name:       "perawat test kelurahan " + uuid.NewString(),
			Profesion:  "nurse",
		}

		res, err := srv.StoreData(data, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)
	})

	t.Run("show list", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "nurse", "true")
		assert.Equal(t, err, nil)
		log.Println(res)
	})
}

func TestShowList(t *testing.T) {
	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	t.Run("show list false", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "nurse", "false")
		assert.Equal(t, err, nil)
		log.Println(res)
	})

	t.Run("show list true", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "nurse", "true")
		assert.Equal(t, err, nil)
		log.Println(res)
	})

	t.Run("show list false", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "doctor", "false")
		assert.Equal(t, err, nil)
		log.Println(res)
	})

	t.Run("show list true", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "doctor", "true")
		assert.Equal(t, err, nil)
		log.Println(res)
	})
}

func TestDoctor(t *testing.T) {
	// 6357f6aa495cc77b8bd6e50c

	repo, err := MongoDeclareForTesting()
	if err != nil {
		t.Fatal(err)
	}

	srv := New(repo)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	t.Run("store doctor and delete", func(t *testing.T) {
		data := &entities.Officer{
			DistrictID: "6319e9dcb0aa169ea947cf68",
			Name:       "dokter test kelurahan " + uuid.NewString(),
			Profesion:  "doctor",
		}

		res, err := srv.StoreData(data, &auth.Token{})
		assert.Equal(t, err, nil)
		log.Println(res)

		delData := res.(*entities.Officer)

		err = srv.SoftDelete(&auth.Token{}, utils.ObjectIDHex(delData.ID))
		assert.Equal(t, err, nil)
	})

	t.Run("fail soft test delete", func(t *testing.T) {
		err := srv.SoftDelete(&auth.Token{}, "6357f6aa495cc77b8bd6e50c")
		assert.NotEqual(t, err, nil)
	})

	t.Run("show list", func(t *testing.T) {
		res, err := srv.GetList("6319e9dcb0aa169ea947cf68", "doctor", "true")
		assert.Equal(t, err, nil)
		log.Println(res)
	})
}

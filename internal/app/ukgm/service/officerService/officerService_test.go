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

/*
	Cipete Selatan
	Gandaria Selatan
	Cilandak Barat
	Lebak Bulus
	Pondok Labu
	test kelurahan
*/
var districtIdList = []string{"6319e9dbb0aa169ea947cf63", "6319e9dbb0aa169ea947cf64", "6319e9dbb0aa169ea947cf65",
	"6319e9dbb0aa169ea947cf66",
	"6319e9dcb0aa169ea947cf67", "6319e9dcb0aa169ea947cf68"}

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

	srv := New(repo)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println(len(districtIdList), len(doctoNameDistrictList), len(nurseNameDistrictList), len(nurseNameKecamatanList), len(doctorNameKecamatanList))

	// input dokter kelurahan
	for i := 0; i < len(districtIdList); i++ {
		officerData := &entities.Officer{
			DistrictID: districtIdList[i],
			Name:       doctoNameDistrictList[i],
			Level:      "kelurahan",
			Profesion:  "doctor",
		}

		log.Println(officerData)
		log.Println(srv)
		res, err := srv.StoreData(officerData, &auth.Token{})
		log.Println(err)
		log.Println(res)
		assert.Equal(t, err, nil)
		log.Println(err)
		log.Println(res)
	}
}

package districtservice

import (
	"context"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	district_repo "ukgm/internal/app/ukgm/mongoRepository/districtRepo"
	patient_repo "ukgm/internal/app/ukgm/mongoRepository/patientRepo"
	unit_repo "ukgm/internal/app/ukgm/mongoRepository/unitRepo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DistrictService struct {
	Repo        *district_repo.DistrictRepo
	UnitRepo    *unit_repo.UnitRepo
	PatientRepo *patient_repo.PatientRepo
}

const DistrictCollection = "districts"

func New(repo *mongorepository.MongoRepo) *DistrictService {
	return &DistrictService{
		Repo:        &district_repo.DistrictRepo{MongoRepo: repo},
		UnitRepo:    &unit_repo.UnitRepo{MongoRepo: repo},
		PatientRepo: &patient_repo.PatientRepo{MongoRepo: repo},
	}
}

func (d *DistrictService) ShowDistrict(filter bson.M, opt *options.FindOneOptions) (result *entities.District, err error) {
	return d.Repo.ShowDistrict(filter, opt, context.TODO())
}

func (d *DistrictService) FetchDistrictList(filter bson.M, opt *options.FindOptions) (result []entities.DistrictList, err error) {
	return d.Repo.FetchDistrictList(filter, opt, context.TODO())
}

func (d *DistrictService) FindManyDistricts(filter bson.M, opt *options.FindOptions) (result []entities.District, err error) {
	return d.Repo.FetchManyDistricts(filter, opt, context.TODO())
}

package doctor_service

import (
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	doctor_repo "ukgm/internal/app/ukgm/mongoRepository/doctorRepo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DoctorService struct {
	Repo *doctor_repo.DoctorRepo
}

const DoctorCollection = "doctors"

func New(Repo *mongorepository.MongoRepo) (result *DoctorService) {
	result = &DoctorService{
		Repo: &doctor_repo.DoctorRepo{MongoRepo: Repo},
	}
	return
}

func (doctor *DoctorService) StoreDoctor(data *entities.Doctor, auth *entities.AuthUser) (result *entities.Doctor, err error) {
	data.Timestamp.CreateNew(auth)
	return doctor.Repo.StoreDoctor(data)
}

func (doctor *DoctorService) ShowDoctor(filter bson.M, opttion *options.FindOneOptions) (result *entities.Doctor, err error) {
	return doctor.Repo.ShowDoctor(filter, opttion)
}

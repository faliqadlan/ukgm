package doctor_service

import (
	"errors"
	"ukgm/internal/app/ukgm/model/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (doctor *DoctorService) PrepareStoreDoctor(data *entities.Doctor, auth *entities.AuthUser) (*entities.Doctor, error) {
	return doctor.StoreDoctor(data, auth)
}

func (doctor *DoctorService) ShowDoctorbyID(ID string) (result *entities.Doctor, err error) {
	oID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("Failed to parse rip id")
	}

	var filter = bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oID,
	}

	result, err = doctor.ShowDoctor(filter, nil)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("Failed to parse rip id")
	}

	return result, nil
}

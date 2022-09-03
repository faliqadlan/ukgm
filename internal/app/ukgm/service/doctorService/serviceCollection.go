package doctor_service

import "ukgm/internal/app/ukgm/model/entities"

type DoctorServiceColl interface {
	PrepareStoreDoctor(data *entities.Doctor, auth *entities.AuthUser) (*entities.Doctor, error)
	ShowDoctorbyID(ID string) (result *entities.Doctor, err error)
}

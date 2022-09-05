package doctor_service

import (
	"ukgm/internal/app/ukgm/model/entities"

	"firebase.google.com/go/auth"
)

type DoctorServiceColl interface {
	PrepareStoreDoctor(data *entities.Doctor, auth *auth.Token) (*entities.Doctor, error)
	ShowDoctorbyID(ID string) (result *entities.Doctor, err error)
}

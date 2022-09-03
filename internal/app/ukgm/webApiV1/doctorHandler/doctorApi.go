package doctor_handler

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

type DoctorApi struct {
	srv *service.ServiceNew
}

func NewDoctorApi(srv *service.ServiceNew) *DoctorApi {
	return &DoctorApi{
		srv: srv,
	}
}

func (api *DoctorApi) Init(r chi.Router) {
	r.Post("/", api.StoreDoctor)
}

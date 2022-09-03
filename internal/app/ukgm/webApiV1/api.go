package webapiv1

import (
	"ukgm/internal/app/ukgm/service"
	doctor_handler "ukgm/internal/app/ukgm/webApiV1/doctorHandler"
	simple "ukgm/internal/app/ukgm/webApiV1/simpleHandler"

	"github.com/go-chi/chi"
)

type Api struct {
	srv       *service.ServiceNew
	simple    *simple.SimpleApi
	doctorApi *doctor_handler.DoctorApi
}

func NewApi(service *service.ServiceNew, simple *simple.SimpleApi) (result *Api) {
	result = &Api{
		srv:       service,
		simple:    simple,
		doctorApi: doctor_handler.NewDoctorApi(service),
	}
	return
}

func (api *Api) Init(r chi.Router) {
	r.Route("/simple", func(r chi.Router) {
		r.Route("/", api.simple.Init)
	})
	r.Route("/doctor", func(r chi.Router) {
		r.Route("/", api.doctorApi.Init)
	})
}

package webapiv1

import (
	"ukgm/internal/app/ukgm/service"
	district_handler "ukgm/internal/app/ukgm/webApiV1/districtHandler"
	doctor_handler "ukgm/internal/app/ukgm/webApiV1/doctorHandler"
	login "ukgm/internal/app/ukgm/webApiV1/loginHandler"
	simple "ukgm/internal/app/ukgm/webApiV1/simpleHandler"

	"github.com/go-chi/chi"
)

type Api struct {
	srv         *service.ServiceNew
	simple      *simple.SimpleApi
	doctorApi   *doctor_handler.DoctorApi
	loginApi    *login.LoginApi
	districtApi *district_handler.DistrictApi
}

func NewApi(service *service.ServiceNew, loginApi *login.LoginApi) (result *Api) {
	result = &Api{
		srv:         service,
		simple:      simple.NewSimpleApi(service),
		doctorApi:   doctor_handler.NewDoctorApi(service),
		districtApi: district_handler.NewDistrictApi(service),
		loginApi:    loginApi,
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
	r.Route("/login", func(r chi.Router) {
		r.Route("/", api.loginApi.Init)
	})

	r.Route("/district", func(r chi.Router) {
		r.Route("/", api.districtApi.Init)
	})
}

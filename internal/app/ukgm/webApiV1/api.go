package webapiv1

import (
	"ukgm/internal/app/ukgm/service"
	diagnose_handler "ukgm/internal/app/ukgm/webApiV1/diagnoseHandler"
	district_handler "ukgm/internal/app/ukgm/webApiV1/districtHandler"
	doctor_handler "ukgm/internal/app/ukgm/webApiV1/doctorHandler"
	login "ukgm/internal/app/ukgm/webApiV1/loginHandler"
	patient_handler "ukgm/internal/app/ukgm/webApiV1/patientHandler"
	simple "ukgm/internal/app/ukgm/webApiV1/simpleHandler"
	unit_handler "ukgm/internal/app/ukgm/webApiV1/unitHandler"
	user_handler "ukgm/internal/app/ukgm/webApiV1/userHandler"

	"github.com/go-chi/chi"
)

type Api struct {
	srv         *service.ServiceNew
	simple      *simple.SimpleApi
	doctorApi   *doctor_handler.DoctorApi
	loginApi    *login.LoginApi
	districtApi *district_handler.DistrictApi
	userApi     *user_handler.UserApi
	unitApi     *unit_handler.UnitApi
	patientApi  *patient_handler.PatientApi
	diagnoseApi *diagnose_handler.DiagnoseApi
}

func NewApi(service *service.ServiceNew, loginApi *login.LoginApi) (result *Api) {
	result = &Api{
		srv:         service,
		simple:      simple.NewSimpleApi(service),
		doctorApi:   doctor_handler.NewDoctorApi(service),
		districtApi: district_handler.NewDistrictApi(service),
		userApi:     user_handler.NewUserApi(service),
		loginApi:    loginApi,
		unitApi:     unit_handler.NewUnitApi(service),
		patientApi:  patient_handler.NewPatientApi(service),
		diagnoseApi: diagnose_handler.NewDiagnoseApi(service),
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

	r.Route("/user", func(r chi.Router) {
		r.Route("/", api.userApi.Init)
	})

	r.Route("/unit", func(r chi.Router) {
		r.Route("/", api.unitApi.Init)
	})
	r.Route("/patient", func(r chi.Router) {
		r.Route("/", api.patientApi.Init)
	})

	r.Route("/diagnose", func(r chi.Router) {
		r.Route("/", api.diagnoseApi.Init)
	})
}

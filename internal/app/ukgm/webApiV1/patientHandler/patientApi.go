package patient_handler

import (
	"ukgm/internal/app/ukgm/model/request"
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

const (
	PatientUrl = "/api/v1/patient"
)

type PatientApi struct {
	srv *service.ServiceNew
}

func NewPatientApi(srv *service.ServiceNew) *PatientApi {
	return &PatientApi{srv: srv}
}

func (p *PatientApi) Init(r chi.Router) {
	r.With(request.PageRequestNewCtx).Get("/", p.FetchIndex)
	r.Post("/", p.StoreData)
}

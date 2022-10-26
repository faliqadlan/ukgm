package officer_handler

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

const (
	DiagnoseUrl = "/api/v1/officer"
)

type OfficerApi struct {
	srv *service.ServiceNew
}

func NewOfficerApi(srv *service.ServiceNew) *OfficerApi {
	return &OfficerApi{srv: srv}
}

func (o *OfficerApi) Init(r chi.Router) {

	r.Post("/", o.StoreData)
	r.Get("/", o.GetList)
	r.Route("/{officer_id}", func(r chi.Router) {
		r.Delete("/", o.SoftDelete)
	})
}

package diagnose_handler

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

const (
	DiagnoseUrl = "/api/v1/diagnose"
)

type DiagnoseApi struct {
	srv *service.ServiceNew
}

func NewDiagnoseApi(srv *service.ServiceNew) *DiagnoseApi {
	return &DiagnoseApi{srv: srv}
}

func (d *DiagnoseApi) Init(r chi.Router) {
	r.Get("/", d.GetDetail)
	r.Get("/list", d.GetList)
	r.Post("/", d.StoreData)
}

package district_handler

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

type DistrictApi struct {
	srv *service.ServiceNew
}

func NewDistrictApi(srv *service.ServiceNew) *DistrictApi {
	return &DistrictApi{srv: srv}
}

func (d *DistrictApi) Init(r chi.Router) {
	r.Route("/{district_id}", func(r chi.Router) {
		r.Use(d.DistrictCtx)
		r.Get("/", d.ShowDistrictbyID)
	})
}

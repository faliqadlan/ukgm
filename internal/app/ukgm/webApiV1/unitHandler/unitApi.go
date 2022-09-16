package unit_handler

import (
	"ukgm/internal/app/ukgm/service"

	"github.com/go-chi/chi"
)

type UnitApi struct {
	srv *service.ServiceNew
}

func NewUnitApi(srv *service.ServiceNew) *UnitApi {
	return &UnitApi{srv: srv}
}

func (u *UnitApi) Init(r chi.Router) {
	r.Post("/", u.StoreUnit)
	r.Get("/", u.GetUnitList)
}

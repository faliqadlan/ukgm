package webapiv1

import (
	"github.com/faliqadlan/ukgm/internal/app/ukgm/service"
	"github.com/go-chi/chi"
)

type Api struct {
	srv *service.ServiceNew
}

func NewApi(service *service.ServiceNew) (resut *Api) {
	resut = &Api{srv: service}
	return
}

func (api *Api) Init(r chi.Router) {
	
}
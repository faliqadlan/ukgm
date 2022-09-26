package district_handler

import (
	"context"
	"net/http"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (d *DistrictApi) DistrictCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		districtID := chi.URLParam(r, "district_id")

		data, err := d.srv.ShowDistrictbyID(districtID)
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
			return
		}

		ctx := context.WithValue(r.Context(), request.GetDistrictCtxKey(), data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (d *DistrictApi) ShowDistrictbyID(w http.ResponseWriter, r *http.Request) {
	// get data context middleware
	district := r.Context().Value(request.GetDistrictCtxKey()).(*entities.DistrictResponse)

	responses.SendJSON(responses.NewResponse(district), w)
}

func (d *DistrictApi) GetDistrictList(w http.ResponseWriter, r *http.Request) {
	districtList, err := d.srv.GetDistrictList()
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}
	responses.SendJSON(responses.NewResponse(districtList), w)
}

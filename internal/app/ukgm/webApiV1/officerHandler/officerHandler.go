package officer_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	"ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (o *OfficerApi) StoreData(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	var mainRequest *entities.Officer

	if err := json.NewDecoder(r.Body).Decode(&mainRequest); err != nil {
		zap.S().Error(err)
		response.SendErrorJSON(response.PrintErrorMessage(err), w)
		return
	}

	res, err := o.srv.OfficerSrvStoreData.StoreData(mainRequest, auth)
	if err != nil {
		zap.S().Error(err)
		response.SendErrorJSON(response.PrintErrorMessage(err), w)
		return
	}
	response.SendJSON(response.NewResponse(res), w)
}

func (o *OfficerApi) SoftDelete(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	officerId := chi.URLParam(r, "officer_id")

	err := o.srv.OfficerSrvSoftDelete.SoftDelete(auth, officerId)
	if err != nil {
		zap.S().Error(err)
		response.SendErrorJSON(response.PrintErrorMessage(err), w)
		return
	}
	response.SendJSON(response.NewResponse(fmt.Sprintf("success delete officer %v", officerId)), w)
}

func (o *OfficerApi) GetList(w http.ResponseWriter, r *http.Request) {

	qParam := r.URL.Query()

	districtId := qParam.Get("district_id")
	profesion := qParam.Get("profesion")
	withKecamatan := qParam.Get("with_kecamatan")

	res, err := o.srv.OfficerSrvGetList.GetList(districtId, profesion, withKecamatan)
	if err != nil {
		zap.S().Error(err)
		response.SendErrorJSON(response.PrintErrorMessage(err), w)
		return
	}
	response.SendJSON(response.NewResponse(res), w)
}

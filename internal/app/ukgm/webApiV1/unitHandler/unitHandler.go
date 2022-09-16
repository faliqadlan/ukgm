package unit_handler

import (
	"encoding/json"
	"net/http"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"go.uber.org/zap"
)

func (u *UnitApi) StoreUnit(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	var mainRequest *entities.Unit
	if err := json.NewDecoder(r.Body).Decode(&mainRequest); err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	res, err := u.srv.UnitSrvStoreData.StoreData(mainRequest, auth)

	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	responses.SendJSON(responses.NewResponse(res), w)
}

func (u *UnitApi) GetUnitList(w http.ResponseWriter, r *http.Request) {
	var districtId = r.URL.Query().Get("district_id")

	unitList, err := u.srv.UnitSrvGetList.GetList(districtId)

	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	responses.SendJSON(responses.NewResponse(unitList), w)
}

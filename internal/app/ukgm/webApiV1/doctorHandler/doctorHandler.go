package doctor_handler

import (
	"encoding/json"
	"net/http"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"go.uber.org/zap"
)

func (doctor *DoctorApi) StoreDoctor(w http.ResponseWriter, r *http.Request) {

	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	var mainRequest *entities.Doctor

	if err := json.NewDecoder(r.Body).Decode(&mainRequest); err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}
	var districtId = r.URL.Query().Get("district_id")

	mainRequest.DistrictID = districtId

	res, err := doctor.srv.PrepareStoreDoctor(mainRequest, auth)

	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	responses.SendJSON(responses.NewResponse(res), w)
}

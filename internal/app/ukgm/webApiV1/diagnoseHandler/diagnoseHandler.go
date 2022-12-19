package diagnose_handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"go.uber.org/zap"
)

func (d *DiagnoseApi) DiagnoseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		patientId := r.URL.Query().Get("patient_id")
		log.Println(patientId)
		data, err := d.srv.DiagnoseSrvGetDetail.GetDetail(patientId)
		log.Println(err)
		if err != nil {
			zap.S().Error(err)
			responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
			return
		}

		ctx := context.WithValue(r.Context(), request.GetDiagnoseCtxKey(), data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (d *DiagnoseApi) StoreData(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	var mainRequest *entities.Diagnose

	if err := json.NewDecoder(r.Body).Decode(&mainRequest); err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	query := r.URL.Query()

	var dataInput []interface{}

	switch query.Get("replace") {
	case "true":
		dataInput = []interface{}{mainRequest, true}
	case "false":
		dataInput = []interface{}{mainRequest, false}
	default:
		dataInput = []interface{}{mainRequest, true}
	}

	res, err := d.srv.DiagnoseSrvStoreData.StoreData(dataInput, auth)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}
	responses.SendJSON(responses.NewResponse(res), w)
}

func (d *DiagnoseApi) GetDetail(w http.ResponseWriter, r *http.Request) {

	patientId := r.URL.Query().Get("patient_id")

	data, err := d.srv.DiagnoseSrvGetDetail.GetDetail(patientId)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	responses.SendJSON(responses.NewResponse(data), w)
}

func (d *DiagnoseApi) GetList(w http.ResponseWriter, r *http.Request) {
	districtId := r.URL.Query().Get("district_id")
	unitId := r.URL.Query().Get("unit_id")
	lastMonthsNumber := r.URL.Query().Get("last_months_number")

	res, err := d.srv.DiagnoseSrvGetList.GetList(districtId, unitId, lastMonthsNumber)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}
	responses.SendJSON(responses.NewResponse(res), w)
}

func (d *DiagnoseApi) GetRecord(w http.ResponseWriter, r *http.Request) {
	patientId := r.URL.Query().Get("patient_id")

	res, err := d.srv.DiagnoseSrvGetList.GetList("record", patientId)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}
	responses.SendJSON(responses.NewResponse(res), w)
}
package patient_handler

import (
	"encoding/json"
	"net/http"
	"ukgm/internal/app/ukgm/model/request"
	responses "ukgm/internal/app/ukgm/model/response"

	"firebase.google.com/go/auth"
	"go.uber.org/zap"
)

func (p *PatientApi) StoreData(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE
	auth := r.Context().Value(request.GetAuthCtxKey()).(*auth.Token)

	var mainRequest *request.Patient

	if err := json.NewDecoder(r.Body).Decode(&mainRequest); err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	patReq, err := mainRequest.ToEntity()
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	res, err := p.srv.PatientSrvStoreData.StoreData(patReq, auth)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}
	responses.SendJSON(responses.NewResponse(res), w)
}

func (p *PatientApi) FetchIndex(w http.ResponseWriter, r *http.Request) {
	// GET DATA CONTEXT MIDDLEWARE

	pr := r.Context().Value(request.GetPageRequestCtxKey()).(*request.PageRequest)
	res, count, err := p.srv.PatientSrvFetchIndex.FetchIndex(pr)
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	var resPR *responses.PageRequest
	resPR, err = pr.ToResponsePR()
	if err != nil {
		zap.S().Error(err)
		responses.SendErrorJSON(responses.PrintErrorMessage(err), w)
		return
	}

	responses.SendJSON(responses.SetPaginate(r, resPR, res, count, PatientUrl), w)
}

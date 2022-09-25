package patient_repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/ukgm/model/request"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type PatientRepo struct {
	*mongorepository.MongoRepo
}

const PatientCollection = "patients"

func (p *PatientRepo) StorePatient(data *entities.Patient, ctx context.Context) (result *entities.Patient, err error) {
	mongoResult, err := p.Store(ctx, PatientCollection, data, nil)

	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data.ID = utils.ObjectIDHex(mongoResult.InsertedID)
	return data, nil
}

func (p *PatientRepo) GenerateMRNumber(ctx context.Context) (result string, err error) {
	total, err := p.Count(ctx, PatientCollection, nil)
	if err != nil {
		zap.S().Error(err)
		return "", err
	}
	total++
	mrNumber := strconv.Itoa(total)
	ind := 4 - len(mrNumber)
	for index := 0; index < ind; index++ {
		mrNumber = "0" + mrNumber
	}
	result = "CLST" + mrNumber
	return result, nil
}

func (p *PatientRepo) Showpatient(filter bson.M, options *options.FindOneOptions, ctx context.Context) (result *entities.Patient, err error) {
	result = &entities.Patient{}
	err = p.Show(ctx, PatientCollection, filter, options).Decode(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data, err := json.Marshal(result)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	s := string(data)
	message := fmt.Sprintf("Found a single document : " + s + "\n")
	p.MongoRepo.BasicMongo.LogPrint(message)

	return result, nil
}

func (p *PatientRepo) FetchIndexPatient(ctx context.Context, pr *request.PageRequest) (result []entities.PatientList, count int, err error) {
	result = []entities.PatientList{}
	var patientTemp *entities.Patient
	filter, options := pr.QueryPipeline(patientTemp.GetSearchableField())
	count, err = p.Count(ctx, PatientCollection, filter.Map())
	if count == 0 || err != nil {
		if err != nil {
			zap.S().Error(err)
		}
		return
	}

	getCursor, err := p.Find(ctx, PatientCollection, filter.Map(), options)
	if err != nil {
		zap.S().Error(err)
		return
	}
	for getCursor.Next(ctx) {
		var elem entities.Patient
		if err = getCursor.Decode(&elem); err != nil {
			zap.S().Error(err)
			return
		}
		result = append(result, *elem.ToPatientList())
	}
	return
}

func (p *PatientRepo) UpdatePatient(data *entities.Patient, ctx context.Context) (err error) {
	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update district")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = p.Update(ctx, PatientCollection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

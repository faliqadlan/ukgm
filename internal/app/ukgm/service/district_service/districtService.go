package districtservice

import (
	"context"
	"errors"
	"fmt"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/ukgm/model/entities"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (d *DistrictService) ShowDistrictbyID(id string) (result *entities.DistrictResponse, err error) {

	errGeneral := fmt.Errorf("failed to show district")
	var ctx = context.TODO()
	var session mongo.Session

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("failed to parse id")
	}

	if session, err = d.Repo.Client.StartSession(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		zap.S().Error(err)
		return nil, errGeneral
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		var filter = bson.M{
			"timestamp.deleted_at": nil,
			"_id":                  oID,
		}

		res, err := d.Repo.ShowDistrict(filter, nil, sc)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		result = res.ToDistrictResponse()

		// handle district list
		filter = bson.M{
			"timestamp.deleted_at": nil,
			"district_id":          utils.ObjectIDHex(oID),
		}

		units, err := d.UnitRepo.FecthUnitDistrictList(sc, filter, nil)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		for i := range units {
			filter = bson.M{
				"timestamp.deleted_at": nil,
				"unit_id":              units[i].ID,
			}
			count, err := d.PatientRepo.TotalPatients(sc, filter)
			if err != nil {
				zap.S().Error(err)
				return errGeneral
			}
			units[i].TotalPatients = *count
		}
		result.Units = units

		// handle condition overview
		filter = bson.M{
			"timestamp.deleted_at": nil,
			"district_id":          utils.ObjectIDHex(oID),
		}
		patients, err := d.PatientRepo.FetchPatientList(sc, filter, nil)
		if err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		contitionOverview := entities.ConditionOverview{}

		for _, pat := range patients {
			if pat.LastDiagnose != nil {
				switch pat.LastDiagnose.Condition {
				case ukgm.KondisiBaik:
					contitionOverview.Baik += 1
				case ukgm.KondisiCukupBaik:
					contitionOverview.CukupBaik += 1
				case ukgm.KondisiKurangBaik:
					contitionOverview.KurangBaik += 1
				case ukgm.KondisiBuruk:
					contitionOverview.Buruk += 1
				}
			}
		}
		result.ConditionOverview = contitionOverview

		if err = session.CommitTransaction(sc); err != nil {
			zap.S().Error(err)
			return errGeneral
		}

		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.TODO()); abortErr != nil {
			zap.S().Error(err)
			return nil, errGeneral
		}
		zap.S().Error(err)
		return nil, errGeneral
	}

	return result, nil
}

func (d *DistrictService) GetDistrictList() (result []entities.DistrictList, err error) {
	filter := bson.M{
		"timestamp.deleted_at": nil,
	}
	res, err := d.FetchDistrictList(filter, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	return res, nil
}

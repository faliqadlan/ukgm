package diagnose_repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"ukgm/internal/app/ukgm"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DiagnoseRepo struct {
	*mongorepository.MongoRepo
}

const DiagnoseCollection = "diagnoses"

func (d *DiagnoseRepo) StoreData(data *entities.Diagnose, ctx context.Context) (*entities.Diagnose, error) {
	mongoResult, err := d.Store(ctx, DiagnoseCollection, data, nil)

	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	data.ID = utils.ObjectIDHex(mongoResult.InsertedID)
	return data, nil
}

func (d *DiagnoseRepo) ShowData(filter bson.M, options *options.FindOneOptions, ctx context.Context) (result *entities.Diagnose, err error) {
	result = &entities.Diagnose{}
	err = d.Show(ctx, DiagnoseCollection, filter, options).Decode(result)
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
	d.MongoRepo.BasicMongo.LogPrint(message)

	return result, nil
}

func (d *DiagnoseRepo) UpdateData(data *entities.Diagnose, ctx context.Context) (err error) {
	update := *data
	update.ID = nil

	oId, err := primitive.ObjectIDFromHex(utils.ObjectIDHex(data.ID))
	if err != nil {
		zap.S().Error(err)
		return errors.New("error in update diagnose")
	}

	filter := bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oId,
	}

	_, err = d.Update(ctx, DiagnoseCollection, filter, update, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}

func (d *DiagnoseRepo) GetList(input ...string) (result interface{}, err error) {
	var districtId, unitId, lastMonthNumberS string
	switch {
	case input[0] != "" && len(input) > 0:
		districtId = input[0]
		fallthrough
	case input[1] != "" && len(input) > 1:
		unitId = input[1]
		fallthrough
	case input[2] != "" && len(input) > 2:
		lastMonthNumberS = input[2]
	}

	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location).Local()

	if lastMonthNumberS == "" {
		lastMonthNumberS = "1"
	}

	lastMonthNumber, err := strconv.Atoi(lastMonthNumberS)
	if err != nil {
		return nil, err
	}
	var diagnoseProgessArr []entities.DiagnoseProgress

	for i := 0; i < lastMonthNumber; i++ {
		diagnoseProgress := entities.DiagnoseProgress{}
		minMonth := now.AddDate(0, -i, 0)
		diagnoseProgress.Date = minMonth.Format("2006-01")
		var filter bson.M

		switch {
		case districtId != "":
			filter = bson.M{
				"district_id":          districtId,
				"timestamp.deleted_at": nil,
				"timestamp.created_at": bson.M{
					"$gte": primitive.NewDateTimeFromTime(utils.StartOfMonth(minMonth)),
					"$lte": primitive.NewDateTimeFromTime(utils.EndOfMonth(minMonth)),
				},
			}
		case unitId != "":
			filter = bson.M{
				"unit_id":              unitId,
				"timestamp.deleted_at": nil,
				"timestamp.created_at": bson.M{
					"$gte": primitive.NewDateTimeFromTime(utils.StartOfMonth(minMonth)),
					"$lte": primitive.NewDateTimeFromTime(utils.EndOfMonth(minMonth)),
				},
			}
		case unitId != "" && districtId != "":
			filter = bson.M{
				"district_id":          districtId,
				"unit_id":              unitId,
				"timestamp.deleted_at": nil,
				"timestamp.created_at": bson.M{
					"$gte": primitive.NewDateTimeFromTime(utils.StartOfMonth(minMonth)),
					"$lte": primitive.NewDateTimeFromTime(utils.EndOfMonth(minMonth)),
				},
			}
		default:
			filter = bson.M{
				"timestamp.deleted_at": nil,
				"timestamp.created_at": bson.M{
					"$gte": primitive.NewDateTimeFromTime(utils.StartOfMonth(minMonth)),
					"$lte": primitive.NewDateTimeFromTime(utils.EndOfMonth(minMonth)),
				},
			}
		}

		cursor, err := d.Find(context.TODO(), DiagnoseCollection, filter, nil)
		if err != nil {
			return nil, err
		}

		for cursor.Next(context.TODO()) {
			diagnose := entities.Diagnose{}
			if err := cursor.Decode(&diagnose); err != nil {
				return nil, err
			}

			switch diagnose.Condition {
			case ukgm.KondisiBaik:
				diagnoseProgress.ConditionOverview.Baik += 1
			case ukgm.KondisiCukupBaik:
				diagnoseProgress.ConditionOverview.CukupBaik += 1
			case ukgm.KondisiKurangBaik:
				diagnoseProgress.ConditionOverview.KurangBaik += 1
			case ukgm.KondisiBuruk:
				diagnoseProgress.ConditionOverview.Buruk += 1

			}
		}
		diagnoseProgessArr = append(diagnoseProgessArr, diagnoseProgress)
	}
	result = diagnoseProgessArr

	return
}

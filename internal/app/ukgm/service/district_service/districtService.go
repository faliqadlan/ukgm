package districtservice

import (
	"errors"
	"ukgm/internal/app/ukgm/model/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (d *DistrictService) ShowDistrictbyID(id string) (result *entities.District, err error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("failed to parse id")
	}

	var filter = bson.M{
		"timestamp.deleted_at": nil,
		"_id":                  oID,
	}

	result, err = d.ShowDistrict(filter, nil)
	if err != nil {
		zap.S().Error(err)
		return nil, errors.New("error in show district")
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

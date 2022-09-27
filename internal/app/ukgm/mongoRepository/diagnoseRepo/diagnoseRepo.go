package diagnose_repo

import (
	"context"
	"ukgm/internal/app/ukgm/model/entities"
	mongorepository "ukgm/internal/app/ukgm/mongoRepository"
	"ukgm/internal/app/utils"

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
}

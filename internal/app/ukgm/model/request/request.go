package request

import (
	"fmt"
	"time"
	"ukgm/internal/app/ukgm/model/entities"

	"github.com/iman231/go-help"
	"go.uber.org/zap"
)

type Patient struct {
	UnitID     string `json:"unit_id" bson:"_unit_id,omitempty"`
	DistrictID string `json:"district_id" bson:"_district_id,omitempty"`
	Name       string `json:"name" bson:"name" part:"main" properties:"searchable"`
	Dob        string `json:"dob" bson:"dob"`
	Gender     string `json:"gender" bson:"gender"`
	NoHp       string `json:"no_hp" bson:"no_hp"`
}

func (p *Patient) ToEntity() (*entities.Patient, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	dob, err := time.Parse("2006-01-02", p.Dob)
	if err != nil {
		zap.S().Error(err)
		return nil, fmt.Errorf("invalid time format")
	}
	dob = dob.Local().In(loc)

	if !help.InArray(p.Gender, []string{"P", "L"}) {
		zap.S().Error(err)
		return nil, fmt.Errorf("invalid gender input")
	}

	return &entities.Patient{
		UnitID:     p.UnitID,
		DistrictID: p.DistrictID,
		Name:       p.Name,
		Dob:        &dob,
		Gender:     p.Gender,
		NoHp: p.NoHp,
	}, nil

}

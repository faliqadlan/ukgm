package entities

import (
	"ukgm/internal/app/utils"
)

type (
	District struct {
		ID     interface{}        `json:"id" bson:"_id,omitempty"`
		Name   string             `json:"name" bson:"name"`
		Doctor map[string]*Doctor `json:"Doctor" bson:"Doctor"`
		Nurse  map[string]*Nurse  `json:"Nurse" bson:"Nurse"`

		// include
		Timestamp
	}

	DistrictList struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	DistrictResponse struct {
		ID                interface{}     `json:"_id"`
		Name              string          `json:"name"`
		Units             []*UnitDistrict `json:"units"`
		ConditionOverview `json:"condition_overview"`

		// include
		Timestamp
	}

	UnitDistrict struct {
		ID            interface{} `json:"id" bson:"_id,omitempty"`
		DistrictID    string      `json:"district_id" bson:"district_id,omitempty"`
		Name          string      `json:"name" bson:"name"`
		TotalPatients int         `json:"total_patient" bson:"total_patient"`
	}

	ConditionOverview struct {
		Baik       int `json:"baik"`
		CukupBaik  int `json:"cukup_baik"`
		KurangBaik int `json:"kurang_baik"`
		Buruk      int `json:"buruk"`
	}
)

func (d *District) ToDistrictList() *DistrictList {
	return &DistrictList{
		ID:   utils.ObjectIDHex(d.ID),
		Name: d.Name,
	}
}

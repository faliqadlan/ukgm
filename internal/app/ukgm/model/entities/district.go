package entities

import "ukgm/internal/app/utils"

type (
	District struct {
		ID     interface{}      `json:"id" bson:"_id,omitempty"`
		Name   string           `json:"name" bson:"name"`
		Units  map[string]*Unit `json:"units" bson:"units"`
		Doctor map[string]*Doctor        `json:"Doctor" bson:"Doctor"`
		Nurse  map[string]*Nurse         `json:"Nurse" bson:"Nurse"`

		// include
		Timestamp
	}

	DistrictList struct {
		ID   string `json:"id" bson:"id"`
		Name string `json:"name" bson:"name"`
	}
)

func (d *District) ToDistrictList() *DistrictList {
	return &DistrictList{
		ID:   utils.ObjectIDHex(d.ID),
		Name: d.Name,
	}
}

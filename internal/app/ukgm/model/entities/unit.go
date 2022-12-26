package entities

import "ukgm/internal/app/utils"

type Unit struct {
	ID         interface{} `json:"id" bson:"_id,omitempty"`
	DistrictID string      `json:"district_id" bson:"district_id,omitempty"`
	Name       string      `json:"name" bson:"name"`

	// include
	Timestamp
}

type UnitList struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func (u *Unit) ToUnitList() *UnitList {
	return &UnitList{
		ID:   utils.ObjectIDHex(u.ID),
		Name: u.Name,
	}
}

package entities

import "ukgm/internal/app/utils"

type Officer struct {
	ID           interface{} `json:"id" bson:"_id,omitempty"`
	DistrictID   string      `json:"district_id" bson:"district_id,omitempty"`
	Name         string      `json:"name" bson:"name" part:"main" properties:"searchable"`
	Level        string      `json:"level" bson:"level"`
	Profesion    string      `json:"profesion" bson:"profesion"`
	DistrictName string      `json:"district_name"`

	// include
	Timestamp
}

type OfficerList struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Profesion    string `json:"profesion"`
	DistrictName string `json:"district_name"`
}

func (o *Officer) ToOfficerList() OfficerList {
	return OfficerList{
		Name: o.Name, 
		ID: utils.ObjectIDHex(o.ID),
		Profesion: o.Profesion,
		DistrictName: o.DistrictName,
	}
}
